package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html"
	"os"
	"strings"
	"time"

	"github.com/rsrini7/godiff/utils"
)

func GenerateHtml(filename1, filename2 string, info1, info2 os.FileInfo, msg1, msg2 string, data1, data2 [][]byte, is_error bool) {
	outfmt := OutputFormat{
		name1:     filename1,
		name2:     filename2,
		fileinfo1: info1,
		fileinfo2: info2,
	}

	var span string
	if is_error {
		span = "<span class=\"err\">"
	} else {
		span = "<span class=\"msg\">"
	}

	if msg1 != "" {
		outfmt.buf1.WriteString(span)
		write_html_bytes(&outfmt.buf1, []byte(msg1))
		outfmt.buf1.WriteString("</span><br>")
	} else if data1 != nil && len(data1) > 0 {
		html_preview_file(&outfmt.buf1, data1)
	}

	if msg2 != "" {
		outfmt.buf2.WriteString(span)
		write_html_bytes(&outfmt.buf2, []byte(msg2))
		outfmt.buf2.WriteString("</span><br>")
	} else if data2 != nil && len(data2) > 0 {
		html_preview_file(&outfmt.buf2, data2)
	}

	html_file_table(&outfmt)

	out.WriteString("<tr><td class=\"ttd\">")
	out.Write(outfmt.buf1.Bytes())

	out.WriteString("</td><td class=\"ttd\">")
	out.Write(outfmt.buf2.Bytes())

	out.WriteString("</td></tr>\n")
	out.WriteString("</table><br>\n")

	out_release_lock()
}

func html_file_table_unified(outfmt *OutputFormat) {

	if !outfmt.header_printed {
		out_acquire_lock()
		outfmt.header_printed = true
		out.WriteString("<table class=\"tab\"><tr><td class=\"tth\"><span class=\"hdr\">")
		out.WriteString(html.EscapeString(outfmt.name1))
		out.WriteString("</span>")
		if outfmt.fileinfo1 != nil {
			fmt.Fprintf(out, " <span class=\"inf\">%d %s</span>", outfmt.fileinfo1.Size(), outfmt.fileinfo1.ModTime().Format(time.RFC1123))
		}
		out.WriteString("<br><span class=\"hdr\">")
		out.WriteString(html.EscapeString(outfmt.name2))
		out.WriteString("</span>")
		if outfmt.fileinfo2 != nil {
			fmt.Fprintf(out, " <span class=\"inf\">%d %s</span>", outfmt.fileinfo2.Size(), outfmt.fileinfo2.ModTime().Format(time.RFC1123))
		}
		out.WriteString("</td></tr>")
	}
}

func (chg *DiffChangerUnifiedHtml) diff_lines(ops []DiffOp) {

	html_file_table_unified(chg.OutputFormat)
	chg.buf1.Reset()

	for _, v := range ops {
		switch v.op {
		case DIFF_OP_INSERT:
			write_html_lines_unified(&chg.buf1, "add", "+", chg.file2[v.start2:v.end2], -1, v.start2, chg.lineno_width)

		case DIFF_OP_REMOVE:
			write_html_lines_unified(&chg.buf1, "del", "-", chg.file1[v.start1:v.end1], v.start1, -1, chg.lineno_width)

		case DIFF_OP_MODIFY:
			write_html_lines_unified(&chg.buf1, "del", "-", chg.file1[v.start1:v.end1], v.start1, -1, chg.lineno_width)
			write_html_lines_unified(&chg.buf1, "add", "+", chg.file2[v.start2:v.end2], -1, v.start2, chg.lineno_width)

		default:
			write_html_lines_unified(&chg.buf1, "nop", " ", chg.file1[v.start1:v.end1], v.start1, v.start2, chg.lineno_width)
		}
	}

	out.WriteString("<tr><td class=\"ttd\">")
	out.Write(chg.buf1.Bytes())
	out.WriteString("</td></tr>\n")
}

func html_file_table(outfmt *OutputFormat) {

	if !outfmt.header_printed {
		out_acquire_lock()
		outfmt.header_printed = true
		out.WriteString("<table class=\"tab\"><tr><td class=\"tth\"><span class=\"hdr\">")
		out.WriteString(html.EscapeString(outfmt.name1))
		out.WriteString("</span>")
		if outfmt.fileinfo1 != nil {
			fmt.Fprintf(out, "<br><span class=\"inf\">%d %s</span>", outfmt.fileinfo1.Size(), outfmt.fileinfo1.ModTime().Format(time.RFC1123))
		}
		out.WriteString("</td><td class=\"tth\"><span class=\"hdr\">")
		out.WriteString(html.EscapeString(outfmt.name2))
		out.WriteString("</span>")
		if outfmt.fileinfo2 != nil {
			fmt.Fprintf(out, "<br><span class=\"inf\">%d %s</span>", outfmt.fileinfo2.Size(), outfmt.fileinfo2.ModTime().Format(time.RFC1123))
		}
		out.WriteString("</td></tr>")
	}
}

func (chg *DiffChangerHtml) diff_lines(ops []DiffOp) {

	html_file_table(chg.OutputFormat)

	chg.buf1.Reset()
	chg.buf2.Reset()

	for _, v := range ops {
		switch v.op {
		case DIFF_OP_INSERT:
			write_html_blanks(&chg.buf1, v.end2-v.start2)
			write_html_lines(&chg.buf2, "add", chg.file2[v.start2:v.end2], v.start2, chg.lineno_width)

		case DIFF_OP_REMOVE:
			write_html_lines(&chg.buf1, "del", chg.file1[v.start1:v.end1], v.start1, chg.lineno_width)
			write_html_blanks(&chg.buf2, v.end1-v.start1)

		case DIFF_OP_MODIFY:
			chg.buf1.WriteString("<span class=\"upd\">")
			chg.buf2.WriteString("<span class=\"upd\">")

			start1, start2 := v.start1, v.start2

			for start1 < v.end1 && start2 < v.end2 {

				write_html_lineno(&chg.buf1, start1+1, chg.lineno_width)
				write_html_lineno(&chg.buf2, start2+1, chg.lineno_width)

				if flag_suppress_line_changes {
					write_html_bytes(&chg.buf1, chg.file1[start1])
					write_html_bytes(&chg.buf2, chg.file2[start2])
				} else {
					// report on changes within the line
					line1, line2 := chg.file1[start1], chg.file2[start2]
					pos1, cmp1 := split_runes(line1)
					pos2, cmp2 := split_runes(line2)

					change1, change2 := do_diff(cmp1, cmp2)

					if change1 != nil {
						// perform shift boundaries, to make the changes more readable
						shift_boundaries(cmp1, change1, rune_bouundary_score)
						shift_boundaries(cmp2, change2, rune_bouundary_score)

						write_html_line_change(&chg.buf1, line1, pos1, change1)
						write_html_line_change(&chg.buf2, line2, pos2, change2)

						writeDiffCSVDelta(&chg.diffbuf, line2)
					}
				}

				chg.buf1.WriteByte('\n')
				chg.buf2.WriteByte('\n')
				start1++
				start2++
			}

			chg.buf1.WriteString("</span>")
			chg.buf2.WriteString("</span>")

			if start1 < v.end1 {
				write_html_lines(&chg.buf1, "del", chg.file1[start1:v.end1], start1, chg.lineno_width)
				write_html_blanks(&chg.buf2, v.end1-start1)
			}

			if start2 < v.end2 {
				write_html_blanks(&chg.buf1, v.end2-start2)
				write_html_lines(&chg.buf2, "add", chg.file2[start2:v.end2], start2, chg.lineno_width)
			}

		default:
			n1, n2 := v.end1-v.start1, v.end2-v.start2
			maxn := utils.MaxInt(n1, n2)

			if n1 > 0 {
				write_html_lines(&chg.buf1, "nop", chg.file1[v.start1:v.end1], v.start1, chg.lineno_width)
			}
			if n1 < maxn {
				write_html_blanks(&chg.buf1, maxn-n1)
			}

			if n2 > 0 {
				write_html_lines(&chg.buf2, "nop", chg.file2[v.start2:v.end2], v.start2, chg.lineno_width)
			}
			if n2 < maxn {
				write_html_blanks(&chg.buf2, maxn-n2)
			}
		}
	}

	out.WriteString("<tr><td class=\"ttd\">")
	out.Write(chg.buf1.Bytes())
	out.WriteString("</td><td class=\"ttd\">")
	out.Write(chg.buf2.Bytes())
	out.WriteString("</td></tr>\n")

	writeDiffToCSV(chg.diffbuf.Bytes())
}

func writeDiffCSVDelta(buf *bytes.Buffer, line []byte) {
	buf.Write(line)
	buf.WriteString("\n")
}

func writeDiffToCSV(buf []byte) {
	output_csv_file, err := os.Create(flag_csv_delta)
	if err != nil {
		usage(err.Error())
	}
	defer output_csv_file.Close()

	outCSV := bufio.NewWriter(output_csv_file)
	outCSV.WriteString(strings.Join(csvHeaderData, ","))
	outCSV.WriteString("\n")
	outCSV.Write(buf)
	outCSV.Flush()
}

// Write single line with changes
func write_html_line_change(buf *bytes.Buffer, line []byte, pos []int, change []bool) {

	in_chg := false
	for i, end := 0, len(change); i < end; {
		j, c := i+1, change[i]
		for j < end && change[j] == c {
			j++
		}
		if c && !in_chg {
			buf.WriteString("<span class=\"chg\">")
		} else if !c && in_chg {
			buf.WriteString("</span>")
		}
		write_html_bytes(buf, line[pos[i]:pos[j]])
		i, in_chg = j, c
	}
	if in_chg {
		buf.WriteString("</span>")
	}
}

func write_html_lines_unified(buf *bytes.Buffer, class string, mode string, lines [][]byte, start1, start2, lineno_width int) {
	buf.WriteString("<span class=\"")
	buf.WriteString(class)
	buf.WriteString("\">")
	for _, line := range lines {
		if start1 >= 0 {
			start1++
		}
		if start2 >= 0 {
			start2++
		}
		write_html_lineno_unified(buf, mode, start1, start2, lineno_width)

		write_html_bytes(buf, line)
		buf.WriteByte('\n')
	}
	buf.WriteString("</span>")
}

func write_html_blanks(buf *bytes.Buffer, n int) {
	buf.WriteString("<span class=\"nop\">")
	for n > 0 {
		buf.WriteString("<span class=\"lno\"> </span>\n")
		n--
	}
	buf.WriteString("</span>")
}

func write_html_lineno(buf *bytes.Buffer, lineno, width int) {
	if lineno > 0 {
		fmt.Fprintf(buf, "<span class=\"lno\">%-*d </span>", width, lineno)
	} else {
		buf.WriteString("<span class=\"lno\"> </span>")
	}
}

func write_html_lineno_unified(buf *bytes.Buffer, mode string, lineno1, lineno2, width int) {
	buf.WriteString("<span class=\"lno\">")

	if lineno1 > 0 {
		fmt.Fprintf(buf, "%-*d", width, lineno1)
	} else {
		fmt.Fprintf(buf, "%-*s", width, "")
	}

	if lineno2 > 0 {
		fmt.Fprintf(buf, " %-*d ", width, lineno2)
	} else {
		fmt.Fprintf(buf, " %-*s ", width, "")
	}

	buf.WriteString(mode)
	buf.WriteString(" </span>")
}

func write_html_lines(buf *bytes.Buffer, class string, lines [][]byte, lineno, lineno_width int) {
	buf.WriteString("<span class=\"")
	buf.WriteString(class)
	buf.WriteString("\">")
	for _, line := range lines {
		lineno++
		write_html_lineno(buf, lineno, lineno_width)
		write_html_bytes(buf, line)
		buf.WriteByte('\n')
	}
	buf.WriteString("</span>")
}

func html_preview_file(buf *bytes.Buffer, lines [][]byte) {
	n := utils.MinInt(NUM_PREVIEW_LINES, len(lines))
	w := len(fmt.Sprintf("%d", n))
	buf.WriteString("<span class=\"nop\">")
	for lineno, line := range lines[0:n] {
		write_html_lineno(buf, lineno+1, w)
		write_html_bytes(buf, line)
		buf.WriteByte('\n')
	}
	buf.WriteString("</span></span>")
}

//
// Write bytes to buffer, ready to be output as html,
// replace special chars with html-entities
//
func write_html_bytes(buf *bytes.Buffer, line []byte) {
	var esc string
	lasti := 0
	for i, v := range line {
		switch v {
		case '<':
			esc = html_entity_lt
		case '>':
			esc = html_entity_gt
		case '&':
			esc = html_entity_amp
		case '\'':
			esc = html_entity_squote
		case '"':
			esc = html_entity_dquote
		default:
			continue
		}
		buf.Write(line[lasti:i])
		buf.WriteString(esc)
		lasti = i + 1
	}
	buf.Write(line[lasti:])
}
