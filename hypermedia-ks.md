

██╗  ██╗██╗   ██╗██████╗ ███████╗██████╗ 
██║  ██║╚██╗ ██╔╝██╔══██╗██╔════╝██╔══██╗
███████║ ╚████╔╝ ██████╔╝█████╗  ██████╔╝
██╔══██║  ╚██╔╝  ██╔═══╝ ██╔══╝  ██╔══██╗
██║  ██║   ██║   ██║     ███████╗██║  ██║
╚═╝  ╚═╝   ╚═╝   ╚═╝     ╚══════╝╚═╝  ╚═╝
                                         
███╗   ███╗███████╗██████╗ ██╗ █████╗ 
████╗ ████║██╔════╝██╔══██╗██║██╔══██╗
██╔████╔██║█████╗  ██║  ██║██║███████║
██║╚██╔╝██║██╔══╝  ██║  ██║██║██╔══██║
██║ ╚═╝ ██║███████╗██████╔╝██║██║  ██║
╚═╝     ╚═╝╚══════╝╚═════╝ ╚═╝╚═╝  ╚═╝

    https://hypermedia.systems/
            by
        Carson Gross
        Adam Stepinski
        Deniz Aksimsek


## 1 
Hypermedia: A New Generation

*Hypertexts:
    new forms of writing,
    apearing on computer screens,
    that will branch or perform at the reader's command.
    A hypertext is a non-sequential piece of writing;
    only the computer display makes it practical.*
        Ted Nelson, https://archive.org/details/SelectedPapers1977/page/n7/mode/2up

Hypermedia Control
    A hypermedia control is an element in a hypermedia that describes (or controls) some sort of interaction,
    often with a remote server, by encoding information about that interaction directly and completely within itself.
## 2
The Essence of HTML as a Hypermedia

    Anchor Tags

    Form Tags
## 3
So What Isn't Hypermedia?

```html
<button onclick="fetch('/api/v1/contacts/1')
                .then(response => response.json())
                .then(data => updateUI(data)) ">
    Fetch Contact
</button>
```

```json
{
  "id": 42,
  "email" : "json-example@example.org"
}
```

In particular, the code in updateUI() needs to know about the internal structure and meaning of the data.

It needs to know:

    Exactly how the fields in the JSON data object are structured and named.

    How they relate to one another.

    How to update the local data this new data corresponds with.

    How to render this data to the browser.

    What additional actions/API end points can be called with this data.

Single Page Applications

* more interactive and immersive experience than Web 1.0
* smoothly update elements inline on a page without a dramatic reload of the entire document
* use CSS transitions to create nice visual events
* hook into arbitrary events like mouse movements

## 4
Why Use Hypermedia?

* It is an extremely simple approach to buildling web applications
* It is extremely tolerant of content and API changes.
* It leverages tried and true features of web browsers, such as caching.

Pain points in modern web development:

* Single Page Application infrastructure has become extremely complex, often requiring an entire team to manage.
* JSON API churn - constant changes made to JSON APIs to support application needs - has become a major pain point for many application teams.

Two major reasons hypermedia hasn't made a comeback in web development(Javascript Fatigue):

* the expressiveness of HTML as a hypermedia hasn't changed much, if at all, since HTML 2.0, which was released in the mid 1990s
* the interactivity and expressiveness of HTML has remained frozen,
    the demands of web users have continued to increase,
    calling for more and more interactive web applications.
## 5
A Hypermedia Resurgence?

Multi-Page Applications

Svelte.js - a blend of MPA and SPA

## 6
Hypermedia-Driven Applications

```html
<button hx-get="/contacts/1" hx-target="#contact-ui"> <1>
    Fetch Contact
</button>
```

```html
<details>
  <div>
    Contact: HTML Example
  </div>
  <div>
    <a href="mailto:html-example@example.com">Email</a>
  </div>
</details>
```
More examples at:
https://htmx.org/examples/
## 7
When Should You Use Hypermedia?

- When you dont need a huge amount of reactivity
- Simple CRUD applications

Benefits of building around hypermedia:

- Back button simply works
- deep linking will  just work
- focus on business logic
## 8
When Shouldn't You Use Hypermedia?

- high amounts of interactivity(e.g. spreadsheets)

Presenter's notes:

- my experience with hypermedia driven application using htmx
- data-*
- HyperView

## 9
Questions?

## 10
Extras

- example sites:
https://store.dreamsofcode.io/

- REST
An Architecture coined By Roy Fielding for distributed systems
The "Constraints" of REST
* It is a client-server architecture
* It must be stateless; that is, every request contains all information necessary to respond to that request.
* It must allow for caching
* It must have a uniform interface
    - Identification of resources
    - Manipulation of resources through representations
    - Self-descriptive messages
    - HATEOAS
    ```html
    <html>
        <body>
            <h1>Joe Smith</h1>
            <div>
                <div>Email: joe@exmaple.bar</div>
                <div>Status: Active</div>
            </div>
            <p>
                <a href="/contacts/42/archive">Archive</a>
            </p>
        </body>
    </html>

    ```
    ```json
    {
        "name": "Joe Smith",

    ```
    - HATEOAS & API churn
* It is a layered system
* Optionally, it can allow for Code-On-Demand, that is, scripting.
