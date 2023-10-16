package model

type Story struct {
	Id        uint64 `json:"id"`
	Cast      string `json:"cast"`
	Location  string `json:"location"`
	Plot      string `json:"plot"`
	StoryText string `json:"story_text"`
	Model     string `json:"model"`
}

func CreateStory(cast, location, plot, story_text, model string) error {

	statement := `insert into stories("cast", location, plot, story_text, model) values($1, $2, $3, $4, $5);`

	_, err := db.Query(statement, cast, location, plot, story_text, model)

	return err

}

func GetAllStories() ([]Story, error) {

	stories := []Story{}

	statement := `select id, "cast", location, plot, story_text, model from stories order by id desc;`

	rows, err := db.Query(statement)
	if err != nil {
		return stories, err
	}

	defer rows.Close()

	for rows.Next() {
		var cast, location, plot, story_text, model string
		var id uint64

		err := rows.Scan(&id, &cast, &location, &plot, &story_text, &model)
		if err != nil {
			return stories, err
		}
		story := Story{
			Id:        id,
			Cast:      cast,
			Location:  location,
			Plot:      plot,
			StoryText: story_text,
			Model:     model,
		}

		stories = append(stories, story)
	}
	return stories, err

}

func GetStory(id uint64) (Story, error) {
	statement := `select id, "cast", location, plot, story_text, model from stories where id=$1;`

	story := Story{}
	story.Id = id

	row, err := db.Query(statement, id)
	if err != nil {
		return story, err
	}

	for row.Next() {
		var cast, location, plot, story_text, model string
		err := row.Scan(&cast, &location, &plot, &story_text, &model)
		if err != nil {
			return story, err
		}

		story.Cast = cast
		story.Location = location
		story.Plot = plot
		story.StoryText = story_text
		story.Model = model
	}
	return story, err
}

func Delete(id uint64) error {

	statement := `delete from stories where id=$1;`

	_, err := db.Exec(statement, id)
	return err

}
