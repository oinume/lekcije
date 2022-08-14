package model2

import "fmt"

const teacherURLBase = "https://eikaiwa.dmm.com/teacher/index/%v/"

func (o *Teacher) URL() string {
	return fmt.Sprintf(teacherURLBase, o.ID)
}

func (o *Teacher) IsJapanese() bool {
	return o.CountryID == 392
}
