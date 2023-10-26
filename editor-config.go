package main

type Version struct {
	Full  string `json:"full"`
	Major int    `json:"major,string"`
	Minor int    `json:"minor,string"`
	Dot   int    `json:"dot,string"`
}

type TableEditorConfig struct {
	ShowDataDefinition bool `json:"showDataDefinition"`
	EqualizeColumns    bool `json:"equalizeColumns"`
	DefaultToSnapping  bool `json:"defaultToSnapping"`
	DefaultToRunning   bool `json:"defaultToRunning"`
}

type FieldDefaults struct {
	Height           int    `json:"height"`
	Width            int    `json:"width"`
	Padding          string `json:"padding"`
	RemoveWhiteSpace bool   `json:"removeWhiteSpace,omitempty"`
	UseFieldText     bool   `json:"useFieldText,omitempty"`
}

type Dimensions struct {
	Width  int    `json:"width,string"`
	Height int    `json:"height,string"`
	Unit   string `json:"unit"`
}

type EditorConfig struct {
	Info struct {
		SetupTitle       string  `json:"setupTitle"`
		CurrentVersion   Version `json:"currentVersion,omitempty"`
		CreatedInVersion Version `json:"createdInVersion,omitempty"`
	}
	Options struct {
		RendererVersion string            `json:"rendererVersion"`
		GridPrint       string            `json:"gridPrint"`
		Guides          bool              `json:"guides"`
		ShowGrid        bool              `json:"showGrid,omitempty"`
		Tables          TableEditorConfig `json:"tables"`
		FieldDefaults   FieldDefaults     `json:"fieldDefaults"`
		LivePreview     bool              `json:"livePreview,omitempty"`
		Canvas          struct {
			LimitDimensions struct {
				X bool `json:"x"`
				Y bool `json:"y"`
			} `json:"limitDimensions"`
			Dimensions    Dimensions `json:"dimensions"`
			UseWorksheets bool       `json:"useWorksheets"`
			UsePages      bool       `json:"usePages,omitempty"`
			UseWebPages   bool       `json:"useWebPages,omitempty"`
		} `json:"canvas"`
	} `json:"options"`
}
