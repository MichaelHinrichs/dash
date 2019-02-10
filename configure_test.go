package dash_test

import (
	"path/filepath"
	"testing"

	"github.com/itchio/dash"
	"github.com/itchio/wharf/state"
	"github.com/stretchr/testify/assert"
)

func makeConsumer(t *testing.T) *state.Consumer {
	return &state.Consumer{
		OnMessage: func(lvl string, msg string) {
			t.Helper()
			t.Logf("[%s] %s", lvl, msg)
		},
	}
}

func configureParams(t *testing.T) *dash.ConfigureParams {
	return &dash.ConfigureParams{
		Consumer: makeConsumer(t),
	}
}

func fixParams(t *testing.T) *dash.FixPermissionsParams {
	return &dash.FixPermissionsParams{
		Consumer: makeConsumer(t),
		DryRun:   true,
	}
}

func Test_ConfigureWindows(t *testing.T) {
	root := filepath.Join("testdata", "windows")

	v, err := dash.Configure(root, configureParams(t))
	assert.NoError(t, err, "walks without problems")

	assert.EqualValues(t, 4, len(v.Candidates), "finds all candidates on first walk")

	v32 := *v
	(&v32).FilterPlatform("windows", "386")

	assert.EqualValues(t, 1, len(v32.Candidates), "only one candidate left after filtering")
	assert.EqualValues(t, "launcher.bat", v32.Candidates[0].Path, "batch won")

	v64 := *v
	(&v64).FilterPlatform("windows", "amd64")

	assert.EqualValues(t, 1, len(v64.Candidates), "only one candidate left after filtering")
	assert.EqualValues(t, "launcher.bat", v64.Candidates[0].Path, "batch won")
}

func Test_ConfigureWindowsIL2CPP(t *testing.T) {
	root := filepath.Join("testdata", "windows-il2cpp")

	v, err := dash.Configure(root, configureParams(t))
	assert.NoError(t, err, "walks without problems")

	assert.EqualValues(t, 3, len(v.Candidates), "finds all candidates on first walk")

	v32 := *v
	(&v32).FilterPlatform("windows", "386")

	assert.EqualValues(t, 1, len(v32.Candidates), "only one candidate left after filtering")
	assert.EqualValues(t, "game.exe", v32.Candidates[0].Path, "game won")

	v64 := *v
	(&v64).FilterPlatform("windows", "amd64")

	assert.EqualValues(t, 1, len(v64.Candidates), "only one candidate left after filtering")
	assert.EqualValues(t, "game.exe", v64.Candidates[0].Path, "game won")
}

func Test_ConfigureWindowsHtml(t *testing.T) {
	root := filepath.Join("testdata", "windows-html")

	v, err := dash.Configure(root, configureParams(t))
	assert.NoError(t, err, "walks without problems")

	assert.EqualValues(t, 2, len(v.Candidates), "finds all candidates on first walk")

	v32 := *v
	(&v32).FilterPlatform("windows", "386")

	assert.EqualValues(t, 1, len(v32.Candidates), "only one candidate left after filtering")
	assert.EqualValues(t, "game.exe", v32.Candidates[0].Path, "batch won")

	v64 := *v
	(&v64).FilterPlatform("windows", "amd64")

	assert.EqualValues(t, 1, len(v64.Candidates), "only one candidate left after filtering")
	assert.EqualValues(t, "game.exe", v64.Candidates[0].Path, "batch won")
}

func Test_ConfigureDarwin(t *testing.T) {
	root := filepath.Join("testdata", "darwin")

	v, err := dash.Configure(root, configureParams(t))
	assert.NoError(t, err, "walks without problems")
	assert.EqualValues(t, 4, len(v.Candidates), "finds all candidates on first walk")

	fixed, err := dash.FixPermissions(v, fixParams(t))
	assert.NoError(t, err, "fixes permissions without problems")
	assert.EqualValues(t, 3, len(fixed), "had to fix some files")

	vcopy := *v
	(&vcopy).FilterPlatform("darwin", "amd64")

	assert.EqualValues(t, 1, len(vcopy.Candidates), "only one candidate left after filtering")
	assert.EqualValues(t, "Some Grand Game.app", vcopy.Candidates[0].Path, "app wins")
}

func Test_ConfigureDarwinNested(t *testing.T) {
	root := filepath.Join("testdata", "darwin-nested")

	v, err := dash.Configure(root, configureParams(t))
	assert.NoError(t, err, "walks without problems")
	assert.EqualValues(t, 4, len(v.Candidates), "finds all candidates on first walk")

	_, err = dash.FixPermissions(v, fixParams(t))
	assert.NoError(t, err, "fixes permissions without problems")

	vcopy := *v
	(&vcopy).FilterPlatform("darwin", "amd64")

	assert.EqualValues(t, 1, len(vcopy.Candidates), "only one candidate left after filtering")
	assert.EqualValues(t, "osx64/dragonjousting.app", vcopy.Candidates[0].Path, "app wins")
}

func Test_ConfigureDarwinGhost(t *testing.T) {
	root := filepath.Join("testdata", "darwin-ghost")

	v, err := dash.Configure(root, configureParams(t))
	assert.NoError(t, err, "walks without problems")
	assert.EqualValues(t, 3, len(v.Candidates), "finds both execs and one valid app bundle")

	_, err = dash.FixPermissions(v, fixParams(t))
	assert.NoError(t, err, "fixes permissions without problems")

	vcopy := *v
	(&vcopy).FilterPlatform("darwin", "amd64")

	assert.EqualValues(t, 1, len(vcopy.Candidates), "only one candidate left after filtering")
	assert.EqualValues(t, "Awesome Stuff.app", vcopy.Candidates[0].Path, "valid app bundle wins")
}

func Test_ConfigureDarwinSymlink(t *testing.T) {
	root := filepath.Join("testdata", "darwin-symlink")

	v, err := dash.Configure(root, configureParams(t))
	assert.NoError(t, err, "walks without problems")
	assert.EqualValues(t, 1, len(v.Candidates), "finds all candidates on first walk")

	_, err = dash.FixPermissions(v, fixParams(t))
	assert.NoError(t, err, "fixes permissions without problems")

	vcopy := *v
	(&vcopy).FilterPlatform("darwin", "amd64")

	assert.EqualValues(t, 1, len(vcopy.Candidates), "only one candidate left after filtering")
	assert.EqualValues(t, "hello.app", vcopy.Candidates[0].Path, "app wins")
}

func Test_ConfigureLinux(t *testing.T) {
	root := filepath.Join("testdata", "linux")

	v, err := dash.Configure(root, configureParams(t))
	assert.NoError(t, err, "walks without problems")
	assert.EqualValues(t, 5, len(v.Candidates), "finds all candidates on first walk")

	fixed, err := dash.FixPermissions(v, fixParams(t))
	assert.NoError(t, err, "fixes permissions without problems")
	assert.EqualValues(t, 5, len(fixed), "fixed some files")

	vcopy := *v
	(&vcopy).FilterPlatform("linux", "amd64")

	assert.EqualValues(t, 1, len(vcopy.Candidates), "only one candidate left after filtering")
	assert.EqualValues(t, "OpenHexagon", vcopy.Candidates[0].Path, "launcher script wins")
}

func Test_ConfigureLinuxLibs(t *testing.T) {
	root := filepath.Join("testdata", "linux-libs")

	v, err := dash.Configure(root, configureParams(t))
	assert.NoError(t, err, "walks without problems")
	assert.EqualValues(t, 1, len(v.Candidates), "finds all candidates on first walk")

	fixed, err := dash.FixPermissions(v, fixParams(t))
	assert.NoError(t, err, "fixes permissions without problems")
	assert.EqualValues(t, 1, len(fixed), "fixed some files")

	vcopy := *v
	(&vcopy).FilterPlatform("linux", "amd64")

	assert.EqualValues(t, 1, len(vcopy.Candidates), "only one candidate left after filtering")
	assert.EqualValues(t, "game", vcopy.Candidates[0].Path, "binary wins")
}

func Test_ConfigureLinuxDualArch(t *testing.T) {
	root := filepath.Join("testdata", "linux-dual-arch")

	v, err := dash.Configure(root, configureParams(t))
	assert.NoError(t, err, "walks without problems")
	assert.EqualValues(t, 2, len(v.Candidates), "finds all candidates on first walk")

	fixed, err := dash.FixPermissions(v, fixParams(t))
	assert.NoError(t, err, "fixes permissions without problems")
	assert.EqualValues(t, 2, len(fixed), "fixed some files")

	v32 := *v
	(&v32).FilterPlatform("linux", "386")

	assert.EqualValues(t, 1, len(v32.Candidates), "only one candidate left after filtering")
	assert.EqualValues(t, "Game.x86", v32.Candidates[0].Path, "launcher script wins")

	v64 := *v
	(&v64).FilterPlatform("linux", "amd64")

	assert.EqualValues(t, 1, len(v64.Candidates), "only one candidate left after filtering")
	assert.EqualValues(t, "Game.x86_64", v64.Candidates[0].Path, "launcher script wins")
}

func Test_ConfigureLinuxJarFallback(t *testing.T) {
	root := filepath.Join("testdata", "linux-jar-fallback")

	v, err := dash.Configure(root, configureParams(t))
	assert.NoError(t, err, "walks without problems")
	assert.EqualValues(t, 2, len(v.Candidates), "finds all candidates on first walk")

	fixed, err := dash.FixPermissions(v, fixParams(t))
	assert.NoError(t, err, "fixes permissions without problems")
	assert.EqualValues(t, 1, len(fixed), "fixed some files")

	v32 := *v
	(&v32).FilterPlatform("linux", "386")

	assert.EqualValues(t, 1, len(v32.Candidates), "only one candidate left after filtering")
	assert.EqualValues(t, "binary", v32.Candidates[0].Path, "launcher script wins")

	v64 := *v
	(&v64).FilterPlatform("linux", "amd64")

	assert.EqualValues(t, 1, len(v64.Candidates), "only one candidate left after filtering")
	assert.EqualValues(t, "hiddenjar.dat", v64.Candidates[0].Path, "launcher script wins")
}

func Test_ConfigureHtmlMany(t *testing.T) {
	root := filepath.Join("testdata", "html", "many")

	v, err := dash.Configure(root, configureParams(t))
	assert.NoError(t, err, "walks without problems")
	assert.EqualValues(t, 1, len(v.Candidates), "finds all candidates on first walk")

	_, err = dash.FixPermissions(v, fixParams(t))
	assert.NoError(t, err, "fixes permissions without problems")

	vcopy := *v
	(&vcopy).FilterPlatform("darwin", "amd64")

	assert.EqualValues(t, 1, len(vcopy.Candidates), "only one candidate left after filtering")
	assert.EqualValues(t, "index.html", vcopy.Candidates[0].Path, "lowest won")
}

func Test_ConfigureHtmlNested(t *testing.T) {
	root := filepath.Join("testdata", "html", "nested")

	v, err := dash.Configure(root, configureParams(t))
	assert.NoError(t, err, "walks without problems")
	assert.EqualValues(t, 1, len(v.Candidates), "finds all candidates on first walk")

	_, err = dash.FixPermissions(v, fixParams(t))
	assert.NoError(t, err, "fixes permissions without problems")

	vcopy := *v
	(&vcopy).FilterPlatform("darwin", "amd64")

	assert.EqualValues(t, 1, len(vcopy.Candidates), "only one candidate left after filtering")
	assert.EqualValues(t, "ThisContainsStuff/index.html", vcopy.Candidates[0].Path, "lowest won")
}

func Test_ConfigureBiggerIsBetter(t *testing.T) {
	root := filepath.Join("testdata", "bigger-is-better")

	v, err := dash.Configure(root, configureParams(t))
	assert.NoError(t, err, "walks without problems")
	assert.EqualValues(t, 3, len(v.Candidates), "finds all candidates on first walk")

	vcopy := *v
	(&vcopy).FilterPlatform("windows", "amd64")

	assert.EqualValues(t, 3, len(vcopy.Candidates), "three candidates left after filtering")
	assert.EqualValues(t, "tiled.exe", vcopy.Candidates[0].Path, "biggest wins")
}

func Test_ConfigureBlacklist(t *testing.T) {
	root := filepath.Join("testdata", "linux-nodewebkit")

	v, err := dash.Configure(root, configureParams(t))
	assert.NoError(t, err, "walks without problems")
	assert.EqualValues(t, 3, len(v.Candidates), "finds all candidates on first walk")

	vcopy := *v
	(&vcopy).FilterPlatform("linux", "amd64")

	assert.EqualValues(t, 3, len(vcopy.Candidates), "three candidates left after filtering")
	assert.EqualValues(t, "nw", vcopy.Candidates[0].Path, "non-nacl helper wins")
}
