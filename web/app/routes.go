package main

/*
Attach endpoints to handlerfunctions
*/
func (a application) Routes() {

	app.router.HandleFunc("/create-recipe", app.createRecipePage)
	app.router.HandleFunc("/browse", app.allRecipePage)
	app.router.HandleFunc("/myrecipes", app.readRecipePage)
	app.router.HandleFunc("/update-recipe{_id}/", app.updateRecipePage)
	app.router.HandleFunc("/remove-recipe{_id}", app.deleteRecipeHandlerFunc)
	app.router.HandleFunc("/support", app.supportPage)
	app.router.HandleFunc("/mypantry", app.pantryHandler)
	app.router.HandleFunc("/login", app.loginHandler)
	app.router.HandleFunc("/register", app.registerHandler)
	app.router.HandleFunc("/", app.showHandler)
}
