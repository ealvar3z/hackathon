"""Async PyTAK runtime wrapper for ForgeNet."""

from __future__ import annotations

import asyncio
from configparser import ConfigParser
from contextlib import suppress

import pytak


class PyTAKRuntime:
    """Thin lifecycle wrapper around PyTAK TX/RX workers."""

    def __init__(self, cot_url: str, *, debug: bool = False) -> None:
        parser = ConfigParser()
        parser["forgenet"] = {"COT_URL": cot_url, "PYTAK_NO_HELLO": "1"}
        if debug:
            parser["forgenet"]["DEBUG"] = "1"

        self.config = parser["forgenet"]
        self.tx_queue: asyncio.Queue[bytes] = asyncio.Queue()
        self.rx_queue: asyncio.Queue[bytes] = asyncio.Queue()
        self._tasks: list[asyncio.Task[None]] = []

    async def start(self) -> None:
        """Start transport workers."""

        reader, writer = await pytak.protocol_factory(self.config)
        tx_worker = pytak.TXWorker(self.tx_queue, self.config, writer)
        rx_worker = pytak.RXWorker(self.rx_queue, self.config, reader)
        self._tasks = [
            asyncio.create_task(tx_worker.run(), name="forgenet-pytak-tx"),
            asyncio.create_task(rx_worker.run(), name="forgenet-pytak-rx"),
        ]

    async def stop(self) -> None:
        """Stop transport workers."""

        for task in self._tasks:
            task.cancel()
        for task in self._tasks:
            with suppress(asyncio.CancelledError):
                await task
        self._tasks.clear()

    async def publish(self, data: bytes) -> None:
        """Queue a CoT event for transmission."""

        await self.tx_queue.put(data)
        await asyncio.sleep(0)

    async def receive_once(self, timeout: float = 10.0) -> bytes:
        """Wait for a single inbound CoT event."""

        return await asyncio.wait_for(self.rx_queue.get(), timeout=timeout)
