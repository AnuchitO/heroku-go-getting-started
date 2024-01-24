package cycle

type GetCycleProgress struct {
	Name          string `json:"name"`
	PersonalScore int    `json:"personalScore"`
	GoalScore     int    `json:"goalScore"`
	EndDate       string `json:"endDate"`
}
