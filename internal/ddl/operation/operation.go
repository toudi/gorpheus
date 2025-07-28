package operation

const (
	OperationTypeCreateTable = "create-table"
	OperationTypeAddField    = "add-field"
)

type Operation struct {
	CreateTable        *CreateTable        `yaml:"create-table"`
	DropTable          *DropTable          `yaml:"drop-table"`
	AddField           *AddField           `yaml:"add-field"`
	DropField          *DropField          `yaml:"drop-field"`
	AlterField         *AlterField         `yaml:"alter-field"`
	CreateIndex        *CreateIndex        `yaml:"create-index"`
	DropIndex          *DropIndex          `yaml:"drop-index"`
	AddConstraint      *AddConstraint      `yaml:"add-constraint"`
	DropConstraint     *DropConstraint     `yaml:"drop-constraint"`
	CreateView         *CreateView         `yaml:"create-view"`
	DropView           *DropView           `yaml:"drop-view"`
	CreateFunction     *CreateFunction     `yaml:"create-function"`
	DropFunction       *DropFunction       `yaml:"drop-function"`
	RunStatement       *RunStatement       // meant for engine-specific tasks
	RunCode            *RunCode            // meant for "data-like" migration. In other words, a custom code will be executed
	AdjustTableVersion *AdjustTableVersion // internal type which is a consequence of operations unpacking.
}

func (o Operation) Type() string {
	if o.CreateTable != nil {
		return OperationTypeCreateTable
	}
	if o.AddField != nil {
		return OperationTypeAddField
	}

	return ""
}

func (o Operation) TableName() string {
	if o.CreateTable != nil {
		return o.CreateTable.TableName
	}
	if o.AddField != nil {
		return o.AddField.TableName
	}
	if o.AlterField != nil {
		return o.AlterField.TableName
	}
	return ""
}

func (o *Operation) Validate() error {
	if o.AddConstraint != nil {
		return o.AddConstraint.Constraint.Validate()
	}

	return nil
}

type Operations []Operation
