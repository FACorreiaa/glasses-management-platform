# Glasses Management App

To do:

- [] Fix Modals

 I built a platform where the user can insert glasses on  │
│   the platform and then ship them for an user. Now I       │
│   want to change how this works. Now I want that when we   │
│   insert a new pair of glasses on the platform, the user   │
│   should also be shown on the form. So we insert user      │
│   data and the glasses data. Analyse my database deeply    │
│   and my code base. The previous component was using a     │
│   modal to insert a pair of glasses which is causing       │
│   problems because its a client side component and we are using go, templ and htmx.
Build a normal form, make it visually appealing and working on mobile (mobile) where the user can insert user and pair for that user. 
The db schema is already built. The previous query previously inserted only for glasses. And the shipping to an user has another query. Now it should be one query (transaction) where it inserts the user and the pair of glasses to each table. 
The other modal used to insert employees should also be a form now with the same characteristics as the previous one (mobile friendly and good ui).
The project was built with tailwindcss, daisyui and templ with HTMX. 