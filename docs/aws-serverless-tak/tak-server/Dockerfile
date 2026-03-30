FROM public.ecr.aws/docker/library/eclipse-temurin:17-jammy

RUN apt-get update \
    && apt-get install -y net-tools netcat vim certbot curl libxml2-utils unzip python3-pip postgresql ldap-utils dnsutils file cron jq

# Create a non-root user
RUN groupadd --gid 1001 takserver && useradd --create-home --shell /bin/bash --uid 1001 --gid 1001 takserver

ENV HOME=/home/takserver
WORKDIR $HOME

COPY ./ $HOME/

ENV TAK_VERSION=takserver-docker-5.5-RELEASE-45

RUN unzip -o "./takserver-docker-distro/${TAK_VERSION}.zip" \
    && rm -rf ./takserver-docker-distro \
    && rm -rf "./${TAK_VERSION}/docker" \
    && mv ./${TAK_VERSION}/tak ./ \
    && rm -rf "./${TAK_VERSION}"

# Install the AWS CLI
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
    && unzip awscliv2.zip \
    && ./aws/install \
    && rm -rf "awscliv2.zip"

# Create /opt/tak directory and set ownership
RUN mkdir -p /opt/tak && chown -R takserver:takserver /opt/tak

# Create Let's Encrypt directory and set ownership
RUN mkdir -p /etc/letsencrypt && chown -R takserver:takserver /etc/letsencrypt
RUN mkdir -p /var/log/letsencrypt && chown -R takserver:takserver /var/log/letsencrypt
RUN mkdir -p /var/lib/letsencrypt && chown -R takserver:takserver /var/lib/letsencrypt

# Set ownership of all files to the takserver user
RUN chown -R takserver:takserver $HOME

RUN chmod -R +x $HOME/utils
RUN chmod +x ./start.sh

# Start cron service before switching to non-root user
RUN service cron start

# Switch to non-root user
USER takserver

RUN pip3 install --user certbot-dns-route53 boto3

HEALTHCHECK --interval=30s --start-period=60s --retries=5 CMD pgrep -f takserver || exit 1

ENTRYPOINT ["/bin/bash", "-c", "./start.sh"]
