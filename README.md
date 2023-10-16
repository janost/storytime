# Storytime

## What?

Rudimentary self-hosted webapp to generate short stories for children with a self-hosted LLM.

## Why?

Just for playing around and for demonstration purposes.

## How does it work?

Currently it's using ollama hosting mistral-openorca. The model is configurable.

## How do I use it?

Without CUDA:

```
docker compose up
```

With CUDA:
```
docker compose -f docker-compose.cuda.yml up
```

Visit http://127.0.0.1:4000
