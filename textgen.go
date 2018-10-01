package main

import "fmt"

func GenerateText(filename1, filename2 string, msg1, msg2 string) {
	out_acquire_lock()
	if flag_unified_context {
		fmt.Fprintf(out, "<<< %s: %s\n", filename1, msg1)
		fmt.Fprintf(out, ">>> %s: %s\n\n", filename2, msg2)
	} else {
		fmt.Fprintf(out, "--- %s: %s\n", filename1, msg1)
		fmt.Fprintf(out, "+++ %s: %s\n\n", filename2, msg2)
	}
	out_release_lock()
}

func (chg *DiffChangerUnifiedText) diff_lines(ops []DiffOp) {

	if !chg.header_printed {
		out_acquire_lock()
		chg.header_printed = true
		fmt.Fprintf(out, "--- %s\n", chg.name1)
		fmt.Fprintf(out, "+++ %s\n", chg.name2)
	}

	fmt.Fprintf(out, "@@ -%d,%d +%d,%d @@\n", ops[0].start1+1, ops[len(ops)-1].end1-ops[0].start1, ops[0].start2+1, ops[len(ops)-1].end2-ops[0].start2)

	for _, v := range ops {
		switch v.op {
		case DIFF_OP_INSERT, DIFF_OP_REMOVE, DIFF_OP_MODIFY:
			for _, line := range chg.file1[v.start1:v.end1] {
				out.WriteString("- ")
				out.Write(line)
				out.WriteByte('\n')
			}

			for _, line := range chg.file2[v.start2:v.end2] {
				out.WriteString("+ ")
				out.Write(line)
				out.WriteByte('\n')
			}

		default:
			for _, line := range chg.file1[v.start1:v.end1] {
				out.WriteString("  ")
				out.Write(line)
				out.WriteByte('\n')
			}
		}
	}
}

func (chg *DiffChangerText) diff_lines(ops []DiffOp) {

	if !chg.header_printed {
		out_acquire_lock()
		chg.header_printed = true
		fmt.Fprintf(out, "<<< %s\n", chg.name1)
		fmt.Fprintf(out, ">>> %s\n", chg.name2)
	}

	for _, v := range ops {
		switch v.op {
		case DIFF_OP_SAME:
			continue

		case DIFF_OP_INSERT:
			print_line_numbers("a", v.start1-1, -1, v.start2, v.end2)

		case DIFF_OP_REMOVE:
			print_line_numbers("d", v.start1, v.end1, v.start2-1, -1)

		case DIFF_OP_MODIFY:
			print_line_numbers("c", v.start1, v.end1, v.start2, v.end2)
		}

		for _, line := range chg.file1[v.start1:v.end1] {
			out.WriteString("< ")
			out.Write(line)
			out.WriteByte('\n')
		}

		if v.end1 > v.start1 && v.end2 > v.start2 {
			out.WriteString("---\n")
		}

		for _, line := range chg.file2[v.start2:v.end2] {
			out.WriteString("> ")
			out.Write(line)
			out.WriteByte('\n')
		}
	}
}
