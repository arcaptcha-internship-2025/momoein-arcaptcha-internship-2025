## Arcaptcha apartment API

## 🚀 Features

## ⚙️ Setup

1. Create a .env file:

```env

```

2.  Pull the Docker image:

```bash

```

Or manually Build the image:

```bash

```

3. Run the container:

```bash

```

Or using the manually built image:

```bash

```

> 💡 **Note**: If your `.env` file is not in the current directory, provide the full path to it using the `--env-file` flag, like:
>
> ```bash
> docker run -p 8080:8080 --env-file /path/to/.env smaila
> ```

## 📘 API Documentation

After running the container, access the interactive API docs at:
http://127.0.0.1:8080/swagger

## 📤 Example Usage

You can send a POST request to / using curl or tools like Postman:

```bash
curl http://localhost:8080/
```

## 📝 Todo
