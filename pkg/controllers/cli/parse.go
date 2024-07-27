package cli

type ParseResult struct {
	MarkflowFile string
}

func Parse(args []string) ParseResult {
	result := ParseResult{
		MarkflowFile: "./markflow.md",
	}

}
