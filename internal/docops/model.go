package docops

type Kind string

const (
	UnknownKind    Kind = ""
	AsyncAPIKind   Kind = "asyncapi"
	AvroKind       Kind = "avro"
	OpenAPIKind    Kind = "openapi"
	JSONSchemaKind Kind = "jsonschema"
)

type Structure struct {
	Channels   []string
	Operations []string
	Paths      []string
	Properties []string
	Fields     []string
}

type Specification struct {
	Source      string
	Kind        Kind
	SpecVersion string
	ID          string
	Version     string
	Title       string
	Name        string
	Description string
	Structure   Structure
}

type Catalog struct {
	Documents []Specification
}

func NewCatalog(documents []Specification) Catalog {
	return Catalog{Documents: append([]Specification(nil), documents...)}
}

func (c Catalog) All() []Specification {
	return append([]Specification(nil), c.Documents...)
}

func (c Catalog) ByKind(kind Kind) []Specification {
	var documents []Specification
	for _, doc := range c.Documents {
		if doc.Kind == kind {
			documents = append(documents, doc)
		}
	}
	return documents
}

func (c Catalog) Titles(kind Kind) []string {
	var titles []string
	for _, doc := range c.Documents {
		if doc.Kind == kind && doc.Title != "" {
			titles = append(titles, doc.Title)
		}
	}
	return SortedCopy(titles)
}

func (c Catalog) Counts() map[Kind]int {
	counts := make(map[Kind]int)
	for _, doc := range c.Documents {
		counts[doc.Kind]++
	}
	return counts
}
