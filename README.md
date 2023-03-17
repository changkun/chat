# chat

ChatGPT in the command line

## Usage

```bash
$ go install changkun.de/x/chat@latest
$ echo $OPENAI_API_KEY # this should be the OpenAI API key
$ chat
```

## Example

- Use `<Ctrl+D>` to terminate user input. (`Ctrl+Z` + `Enter` in windows)
- Use `<Ctrl+C>` to terminate the chat session.

```
$ chat
> Hi, I'm a chatbot. How can I help you?
> User: How are you? <Ctrl+D>
> Assistant: As an AI language model, I don't have feelings, emotions, or physical experiences. But thank you for asking! How can I assist you today?
> User:
```

## License

MIT
