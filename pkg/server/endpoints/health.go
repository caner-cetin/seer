package endpoints

import "net/http"

// Health is health endpoint which redirects you to a Health song, you should thank this guy for my extremely helpful documentation
//
//	https://github.com/golang/lint/issues/186#issuecomment-171163152
//
// > functions must describe their intents
// idk you will hear this in next two minutes
//
// I wake up in the dark of night (Tell me how does it feel?)
// Do you ever think about me? (Do I need to be?)
// You're sleeping alone tonight (Tonight, right here)
// How does it feel, what do you think?
// Are you happy?
// Faith in everything
// Faith in what you need
// Faith in everything
// Now
// Pain in everything you know
// I wake up with the stars at night
// Tell me what you feel
// And I cannot feel you by me
// It didn't need to be
// You sleep alone tonight (Glory and I hear)
// How does it feel, what do you think?
// Do you still hate me?
func Health(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://www.youtube.com/watch?v=FX0v4sm5dYc", http.StatusFound)
}
