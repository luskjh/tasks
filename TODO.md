- [ ] look for more streamlined way to compose cachec entries into func
      map into registry (currently repeating that logic a lot, but don't
      want to get tightly coupled either)
- [ ] move static assets to /static
- [ ] add public folder in internal/static with favicon etc. for native
      vite handling
- [ ] this may not be the best ergonomics:
      {{ template "shared/vite.html" (vite_entry "js/dashboard.js") }}
      also, we may want to pre-compile these templates
