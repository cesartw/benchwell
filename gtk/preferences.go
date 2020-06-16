package gtk

/*
import (
	"fmt"
	"image/color"

	"bitbucket.org/goreorto/sqlaid/config"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type Preferences struct {
	*gtk.Notebook
}

func (p Preferences) Init() (*Preferences, error) {
	var err error
	p.Notebook, err = gtk.NotebookNew()
	if err != nil {
		return nil, err
	}
	p.Notebook.SetProperty("tab-pos", gtk.POS_LEFT)
	p.Notebook.SetSizeRequest(600, 500)

	uiLabel, uiContent, err := p.SectionGUI()
	if err != nil {
		return nil, err
	}

	p.AppendPage(uiContent, uiLabel)
	p.ShowAll()

	return &p, nil
}

type option struct {
	name       string
	bind       interface{}
	values     []string
	optiontype string
}

func (o option) Widgets() (lbl *gtk.Label, content gtk.IWidget, err error) {
	lbl, err = gtk.LabelNew(o.name)
	if err != nil {
		return nil, nil, err
	}
	switch o.optiontype {
	case "bool":
		sw, err := gtk.SwitchNew()
		if err != nil {
			return nil, nil, err
		}
		sw.SetActive(o.bind.(bool))
		content = sw
	case "multi":
		cb, err := gtk.ComboBoxTextNew()
		if err != nil {
			return nil, nil, err
		}
		for _, v := range o.values {
			cb.Append(v, v)
		}
		cb.SetActiveID(o.bind.(string))
		content = cb
	case "color":
		cp, err := gtk.ColorButtonNew()
		if err != nil {
			return nil, nil, err
		}
		c := color.RGBA{}
		fmt.Sscanf(config.Env.GUI.Editor.Theme.Background, "#%02x%02x%02x", &c.R, &c.G, &c.B)
		cp.ColorChooser.SetRGBA(gdk.NewRGBA(float64(c.R), float64(c.G), float64(c.B)))

	}

	return lbl, content, nil
}

func definition() map[string]interface{} {
	return map[string]interface{}{
		"UI": map[string]interface{}{
			"Dark Mode": option{
				name:       "Dark Mode",
				optiontype: "switch",
				bind:       &config.Env.GUI.DarkMode,
			},
			"Tabs": []option{
				{
					name:       "Connections",
					bind:       &config.Env.GUI.ConnectionTabPosition,
					optiontype: "combo",
					values:     []string{"top", "bottom"},
				},
				{
					name:       "Tables",
					bind:       &config.Env.GUI.TableTabPosition,
					optiontype: "combo",
					values:     []string{"top", "bottom"},
				},
			},
		},
	}
}

func (p *Preferences) Section(name string, section interface{}) (*gtk.Label, gtk.IWidget, error) {
	left := 0
	sectionLbl, err := gtk.LabelNew("<b>" + name + "</b>")
	if err != nil {
		return nil, nil, err
	}
	sectionLbl.SetUseMarkup(true)
	grid, err := gtk.GridNew()
	if err != nil {
		return nil, nil, err
	}
	grid.SetHExpand(true)
	grid.SetVExpand(true)

	var (
		lbl  *gtk.Label
		ctrl gtk.IWidget
	)
	switch t := section.(type) {
	case option:
		lbl, err = gtk.LabelNew(name)
		if err != nil {
			return nil, nil, err
		}
		switch t.optiontype {
		case "bool":
			sw, err := gtk.SwitchNew()
			if err != nil {
				return nil, nil, err
			}
			sw.SetActive(t.bind.(bool))
			ctrl = sw
		}

	case []option:
		lbl, err = gtk.LabelNew("<b>" + name + "</b>")
		if err != nil {
			return nil, nil, err
		}
		lbl.SetUseMarkup(true)
		for _, opt := range t {
			optlbl, err := gtk.LabelNew(name)
			if err != nil {
				return nil, nil, err
			}
			switch opt.optiontype {
			case "bool":
				sw, err := gtk.SwitchNew()
				if err != nil {
					return nil, nil, err
				}
				sw.SetActive(opt.bind.(bool))
				ctrl = sw
			}
		}

	case map[string]interface{}:
	}

	grid.Attach(lbl, left, 0, 1, 1)
	grid.Attach(ctrl, left+1, 0, 1, 1)

	return lbl, ctrl, nil
}

func (p *Preferences) SectionGUI() (*gtk.Label, gtk.IWidget, error) {
	sectionLbl, err := gtk.LabelNew("UI")
	if err != nil {
		return nil, nil, err
	}

	grid, err := gtk.GridNew()
	if err != nil {
		return nil, nil, err
	}
	grid.SetHExpand(true)
	grid.SetVExpand(true)

	tabsLbl, err := gtk.LabelNew("<b>Tabs</b>")
	if err != nil {
		return nil, nil, err
	}
	tabsLbl.SetUseMarkup(true)

	connectionTabLbl, err := gtk.LabelNew("Connection")
	if err != nil {
		return nil, nil, err
	}
	tableTabLbl, err := gtk.LabelNew("Table")
	if err != nil {
		return nil, nil, err
	}
	//TableTabPosition      string `mapstructure:"table_tab_position"`
	//CellWidth             int    `mapstructure:"initial_cell_width"`

	editorLbl, err := gtk.LabelNew("<b>Editor</b>")
	if err != nil {
		return nil, nil, err
	}
	editorLbl.SetUseMarkup(true)

	grid.Attach(tabsLbl, 0, 0, 1, 1)
	grid.Attach(connectionTabLbl, 1, 1, 1, 1)
	grid.Attach(tableTabLbl, 1, 2, 2, 1)
	grid.Attach(editorLbl, 0, 3, 1, 1)

	sectionLbl.Show()
	grid.ShowAll()

	return sectionLbl, grid, nil
}
*/
