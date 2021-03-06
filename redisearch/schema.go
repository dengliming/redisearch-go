package redisearch

// FieldType is an enumeration of field/property types
type FieldType int

// Options are flags passed to the the abstract Index call, which receives them as interface{}, allowing
// for implementation specific options
type Options struct {

	// If set, we will not save the documents contents, just index them, for fetching ids only.
	NoSave bool

	// If set, we avoid saving field bits for each term.
	// This saves memory, but does not allow filtering by specific fields.
	// This is an option that is applied and index level.
	NoFieldFlags bool

	// If set, we avoid saving the term frequencies in the index.
	// This saves memory but does not allow sorting based on the frequencies of a given term within the document.
	// This is an option that is applied and index level.
	NoFrequencies bool

	// If set, , we avoid saving the term offsets for documents.
	// This saves memory but does not allow exact searches or highlighting. Implies NOHL
	// This is an option that is applied and index level.
	NoOffsetVectors bool

	// Set the index with a custom stop-words list, to be ignored during indexing and search time
	// This is an option that is applied and index level.
	// If the list is nil the default stop-words list is used.
	// See https://oss.redislabs.com/redisearch/Stopwords.html#default_stop-word_list
	Stopwords []string
}

// DefaultOptions represents the default options
var DefaultOptions = Options{
	NoSave:          false,
	NoFieldFlags:    false,
	NoFrequencies:   false,
	NoOffsetVectors: false,
	Stopwords:       nil,
}

const (
	// TextField full-text field
	TextField FieldType = iota

	// NumericField numeric range field
	NumericField

	// GeoField geo-indexed point field
	GeoField

	// TagField is a field used for compact indexing of comma separated values
	TagField
)

// Field represents a single field's Schema
type Field struct {
	Name     string
	Type     FieldType
	Sortable bool
	Options  interface{}
}

// TextFieldOptions Options for text fields - weight and stemming enabled/disabled.
type TextFieldOptions struct {
	Weight   float32
	Sortable bool
	NoStem   bool
	NoIndex  bool
}

// TagFieldOptions options for indexing tag fields
type TagFieldOptions struct {
	// Separator is the custom separator between tags. defaults to comma (,)
	Separator byte
	NoIndex   bool
	Sortable  bool
}

// NumericFieldOptions Options for numeric fields
type NumericFieldOptions struct {
	Sortable bool
	NoIndex  bool
}

// NewTextField creates a new text field with the given weight
func NewTextField(name string) Field {
	return Field{
		Name:    name,
		Type:    TextField,
		Options: nil,
	}
}

// NewTextFieldOptions creates a new text field with given options (weight/sortable)
func NewTextFieldOptions(name string, opts TextFieldOptions) Field {
	f := NewTextField(name)
	f.Options = opts
	return f
}

// NewSortableTextField creates a text field with the sortable flag set
func NewSortableTextField(name string, weight float32) Field {
	return NewTextFieldOptions(name, TextFieldOptions{
		Weight:   weight,
		Sortable: true,
	})

}

// NewTagField creates a new text field with default options (separator: ,)
func NewTagField(name string) Field {
	return Field{
		Name:    name,
		Type:    TagField,
		Options: TagFieldOptions{Separator: ',', NoIndex: false},
	}
}

// NewTagFieldOptions creates a new tag field with the given options
func NewTagFieldOptions(name string, opts TagFieldOptions) Field {
	return Field{
		Name:    name,
		Type:    TagField,
		Options: opts,
	}
}

// NewNumericField creates a new numeric field with the given name
func NewNumericField(name string) Field {
	return Field{
		Name: name,
		Type: NumericField,
	}
}

// NewNumericFieldOptions defines a numeric field with additional options
func NewNumericFieldOptions(name string, options NumericFieldOptions) Field {
	f := NewNumericField(name)
	f.Options = options
	return f
}

// NewSortableNumericField creates a new numeric field with the given name and a sortable flag
func NewSortableNumericField(name string) Field {
	f := NewNumericField(name)
	f.Options = NumericFieldOptions{
		Sortable: true,
	}
	return f
}

// Schema represents an index schema Schema, or how the index would
// treat documents sent to it.
type Schema struct {
	Fields  []Field
	Options Options
}

// NewSchema creates a new Schema object
func NewSchema(opts Options) *Schema {
	return &Schema{
		Fields: []Field{},
	}
}

// AddField adds a field to the Schema object
func (m *Schema) AddField(f Field) *Schema {
	if m.Fields == nil {
		m.Fields = []Field{}
	}
	m.Fields = append(m.Fields, f)
	return m
}
