package tabulator

/*

export interface EditorConfig {
  info: {
    setupTitle: string;
    currentVersion?: Version;
    createdInVersion?: Version;
  };
  options: {
    rendererVersion: "v1" | "v2";
    ///Determines if grids are printed in just excel, pdf, both, or neither when in sheets mode
    gridPrint: "never" | "pdf" | "xlsx" | "always";
    guides: boolean;
    showGrid?: boolean;
    tables: TableEditorConfig;
    fieldDefaults: FieldDefaults;
    livePreview?: boolean;
    canvas: {
      limitDimensions: {
        x: boolean;
        y: boolean;
      };
      dimensions: Dimensions;
      useWorksheets: boolean;
      usePages?: boolean;
      useWebPages?: boolean;
    };
  };
}
*/
// export interface Version {
// 	full: string; // 2.54.4,
// 	major: string; //2,
// 	minor: string; //54,
// 	dot: string; //4
//   }

// interface TableEditorConfig {
//   showDataDefinition: boolean;
//   equalizeColumns: boolean;
//   defaultToSnapping: boolean;
//   defaultToRunning: boolean;
// }

// export interface FieldDefaults {
//   height: number;
//   width: number;
//   /** in pixels, i.e. "0px" */
//   padding: string;
//   removeWhiteSpace?: boolean;
//   useFieldText?: boolean;
// }

type Version struct {
	Full  string `json:"full"`
	Major string `json:"major"`
	Minor string `json:"minor"`
	Dot   string `json:"dot"`
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
	Width  int    `json:"width"`
	Height int    `json:"height"`
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
