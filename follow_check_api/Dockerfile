FROM denoland/deno:debian-1.41.0

RUN apt-get update && apt-get install -y --no-install-recommends wget \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . .

# Compile the main app so that it doesn't need to be compiled each startup/entry.
RUN deno cache main.ts

CMD ["run", "--allow-read", "--allow-write", "--allow-net", "--unstable-kv", "--unstable-cron", "main.ts"]
