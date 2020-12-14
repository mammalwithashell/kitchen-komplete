package main

/*
Attach endpoints to handlerfunctions
*/
func (a application) Routes() {

	app.router.HandleFunc("/create-recipe", app.createRecipe)
	app.router.HandleFunc("/browse", app.allRecipePage)
	app.router.HandleFunc("/myrecipes", app.readRecipe)
	app.router.HandleFunc("/recipe/{_id}/", app.recipe)
	app.router.HandleFunc("/edit/{_id}/", app.editRecipe)
	app.router.HandleFunc("/delete/{_id}/", app.recipe)
	app.router.HandleFunc("/support", app.support)
	app.router.HandleFunc("/mypantry", app.pantry)
	app.router.HandleFunc("/login", app.login)
	app.router.HandleFunc("/register", app.register)
	app.router.HandleFunc("/", app.showRecipe)
	app.router.HandleFunc("/profile", app.profile)
	app.router.HandleFunc("/logout", app.logout)
	app.router.HandleFunc("/about", app.about)
	app.router.Use(loggingMiddleware)
}
