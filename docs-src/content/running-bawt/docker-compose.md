---
Title: Docker Compose
Weight: 30
---

The following configuration is how we do local development with Docker Compose.

```yaml
version: '3.7'
services:
  bot:
    image: 'your/docker-image:latest'
    restart: always
    command: 
    - ./bot
    - runserver
    ports:
    - "8080:8080"
    volumes:
    - type: bind  
      source: ./bot.bolt.db
      target: /bot.bolt.db
      read_only: false
    - type: bind  
      source: ./config.yaml
      target: /config.yaml
      read_only: false
```