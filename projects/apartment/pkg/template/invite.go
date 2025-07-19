package template

type InviteData struct {
	Name          string
	EventName     string
	EventDate     string
	EventTime     string
	EventLocation string
	Message       string
	RSVPLink      string
	OrganizerName string
}

// t, err := template.New("invite").ParseFiles("invite_email.html")
// if err != nil {
//     log.Fatal(err)
// }

// data := InviteData{
//     Name:          "Alex",
//     EventName:     "Go Developers Meetup 2025",
//     EventDate:     "August 12, 2025",
//     EventTime:     "6:00 PM",
//     EventLocation: "Tech Hub, Downtown",
//     Message:       "Join us for an evening of networking, talks, and pizza!",
//     RSVPLink:      "https://example.com/rsvp",
//     OrganizerName: "The Go Team",
// }

// var tpl bytes.Buffer
// if err := t.ExecuteTemplate(&tpl, "InviteEmail", data); err != nil {
//     log.Fatal(err)
// }

// Now use tpl.Bytes() or tpl.String() in your email sender
