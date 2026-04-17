package main

import (
	"bytes"
	"flag"
	"image"
	"image/png"
	"io"
	"os"
	"testing"
)

func runMainWithArgs(t *testing.T, args []string) string {
	t.Helper()

	oldArgs := os.Args
	oldStdout := os.Stdout
	oldFlagSet := flag.CommandLine

	defer func() {
		os.Args = oldArgs
		os.Stdout = oldStdout
		flag.CommandLine = oldFlagSet
	}()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	os.Stdout = w
	os.Args = append([]string{"cmd"}, args...)

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	main()

	_ = w.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("failed to read stdout: %v", err)
	}

	return buf.String()
}

func Test_NoParams(t *testing.T) {
	output := runMainWithArgs(t, []string{})

	expected := "Error: At least one value needed for Width or Height\n"
	if output != expected {
		t.Fatalf("expected %q, got %q", expected, output)
	}
}

func Test_NoWidthNoHeight(t *testing.T) {
	output := runMainWithArgs(t, []string{})

	expected := "Error: At least one value needed for Width or Height\n"
	if output != expected {
		t.Fatalf("expected %q, got %q", expected, output)
	}
}

func Test_WidthZero(t *testing.T) {
	output := runMainWithArgs(t, []string{
		"-w", "0",
	})

	expected := "Error: At least one value needed for Width or Height\n"
	if output != expected {
		t.Fatalf("expected %q, got %q", expected, output)
	}
}

func Test_WidthNegative(t *testing.T) {
	output := runMainWithArgs(t, []string{
		"-w", "-100",
	})

	expected := "Error: Negative values not supported for Width/Height\n"
	if output != expected {
		t.Fatalf("expected %q, got %q", expected, output)
	}
}

func Test_HeightZero(t *testing.T) {
	output := runMainWithArgs(t, []string{
		"-h", "0",
	})

	expected := "Error: At least one value needed for Width or Height\n"
	if output != expected {
		t.Fatalf("expected %q, got %q", expected, output)
	}
}

func Test_HeightNegative(t *testing.T) {
	output := runMainWithArgs(t, []string{
		"-h", "-100",
	})

	expected := "Error: Negative values not supported for Width/Height\n"
	if output != expected {
		t.Fatalf("expected %q, got %q", expected, output)
	}
}

func Test_NoFilePath(t *testing.T) {
	output := runMainWithArgs(t, []string{
		"-w", "100",
		"-h", "100",
	})

	expected := "Error: File not Found\n"
	if output != expected {
		t.Fatalf("expected %q, got %q", expected, output)
	}
}

func TestFile_NotImage(t *testing.T) {
	output := runMainWithArgs(t, []string{
		"-f", "test.txt",
		"-w", "100",
		"-h", "100",
	})

	expected := "Error: File not a supported image format\n"
	if output != expected {
		t.Fatalf("expected %q, got %q", expected, output)
	}
}

func TestMain_Success(t *testing.T) {
	dir := t.TempDir()

	inputPath := dir + "/input.png"

	img := image.NewRGBA(image.Rect(0, 0, 10, 20))
	f, err := os.Create(inputPath)
	if err != nil {
		t.Fatal(err)
	}
	err = png.Encode(f, img)
	_ = f.Close()
	if err != nil {
		t.Fatal(err)
	}

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	output := runMainWithArgs(t, []string{
		"-f", inputPath,
		"-w", "5",
		"-h", "5",
	})

	if output != "" {
		t.Fatalf("expected no stdout output, got %q", output)
	}

	outputPath := dir + "/input_resized.png"
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("expected output file to exist: %v", err)
	}

	outFile, err := os.Open(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()

	outImg, _, err := image.Decode(outFile)
	if err != nil {
		t.Fatal(err)
	}

	gotBounds := outImg.Bounds()
	if gotBounds.Dx() != 5 {
		t.Fatalf("expected width 5, got %d", gotBounds.Dx())
	}

	if gotBounds.Dy() != 5 {
		t.Fatalf("expected height 5, got %d", gotBounds.Dy())
	}
}

func TestMain_Success_NoWidth(t *testing.T) {
	dir := t.TempDir()

	inputPath := dir + "/input.png"

	img := image.NewRGBA(image.Rect(0, 0, 10, 20))
	f, err := os.Create(inputPath)
	if err != nil {
		t.Fatal(err)
	}
	err = png.Encode(f, img)
	_ = f.Close()
	if err != nil {
		t.Fatal(err)
	}

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	output := runMainWithArgs(t, []string{
		"-f", inputPath,
		"-h", "5",
	})

	if output != "" {
		t.Fatalf("expected no stdout output, got %q", output)
	}

	outputPath := dir + "/input_resized.png"
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("expected output file to exist: %v", err)
	}

	outFile, err := os.Open(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()

	outImg, _, err := image.Decode(outFile)
	if err != nil {
		t.Fatal(err)
	}

	gotBounds := outImg.Bounds()
	if gotBounds.Dx() != 10 {
		t.Fatalf("expected width 10, got %d", gotBounds.Dx())
	}

	if gotBounds.Dy() != 5 {
		t.Fatalf("expected height 5, got %d", gotBounds.Dy())
	}
}

func TestMain_Success_NoHeight(t *testing.T) {
	dir := t.TempDir()

	inputPath := dir + "/input.png"

	img := image.NewRGBA(image.Rect(0, 0, 10, 20))
	f, err := os.Create(inputPath)
	if err != nil {
		t.Fatal(err)
	}
	err = png.Encode(f, img)
	_ = f.Close()
	if err != nil {
		t.Fatal(err)
	}

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	output := runMainWithArgs(t, []string{
		"-f", inputPath,
		"-w", "5",
	})

	if output != "" {
		t.Fatalf("expected no stdout output, got %q", output)
	}

	outputPath := dir + "/input_resized.png"
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("expected output file to exist: %v", err)
	}

	outFile, err := os.Open(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()

	outImg, _, err := image.Decode(outFile)
	if err != nil {
		t.Fatal(err)
	}

	gotBounds := outImg.Bounds()
	if gotBounds.Dx() != 5 {
		t.Fatalf("expected width 5, got %d", gotBounds.Dx())
	}

	if gotBounds.Dy() != 10 {
		t.Fatalf("expected height 10, got %d", gotBounds.Dy())
	}
}
