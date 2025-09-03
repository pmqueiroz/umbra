```u
io::println(
	html::html(
		html::body(
			html::frag(
				html::div(
					"Hello, World!",
					{ class: "flex items-center" }
				),
				html::div(
					"Click me",
					{
						class: "flex items-center",
						hxPost: "/click",
					}
				),
			),
			{}
		),
		{}
	)
)
```

```
<html>
  <body>
    <div class="flex items-center">
      Hello, World!
    </div>
    <div class="flex items-center" hx-post="/click">
      Click me
    </div>
  </body>
</html>
```
