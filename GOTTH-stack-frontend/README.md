# GOTTH Stack Frontend

- [GOTTH Stack Frontend](#gotth-stack-frontend)
- [Description](#description)
- [Steps](#steps)
- [References](#references)

# Description

Create a fullstack app that uses golang echo, templ, tailwind css, air and htmx.

# Steps

```bash
make deps

# Generate the templ files
templ generate # OR templ generate -watch

# Start the Tailwind build process with --watch
npm run dev

# In another terminal, run the Go server with Air:
air # OR make run

# Then visit http://127.0.0.1:3000/
```

# References
- https://callistaenterprise.se/blogg/teknik/2024/01/08/htmx-with-go-templ/
  - https://github.com/eriklupander/questions-admin-app
- https://templ.guide/quick-start/installation/
- https://github.com/mgechev/revive?tab=readme-ov-file#installation
- https://echo.labstack.com/docs/quick-start
- https://tailwindcss.com/
