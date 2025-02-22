Give me (please) vocabulary language importer where we have a text field that allows us to import a thematic catagoy for the generation of language vocabulary.

When submitting that text field, it should hit an API endpoint (api route in app router) to invoke an LLM chat completions in Groq (LLM) on the server-side and then pass the infomration back to the front-end

It has to create a structured json output like this example:
```
[
  {
    "english": "good",
    "spanish": "bueno"
  },
  {
    "english": "old",
    "spanish": "viejo"
  },
  {
    "english": "busy",
    "spanish": "ocupado"
  }
]
```

The joson that is outputted back to the front-end should be copy-able...so it should be sent to an input field and there shoudl be a copy button so that it can be copied to the clipboard adn then should an alert that it was copied to the user's clipboard.

The app should use app router and the latest version of next.js... and the llm calls should run in an api route on the server-side

