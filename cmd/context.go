package cmd

type context struct {
}

func NewContext() (*context, error) {
	return &context{}, nil
}
