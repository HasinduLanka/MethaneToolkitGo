package methane

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var Excecutables map[string]string = map[string]string{"bash": ""}

func InitExec() error {
	for key := range Excecutables {
		path, err := exec.LookPath(key)
		if err != nil {
			Print("Some excecutables not found. Please make sure you have installed all the needed dependencies.")
			Print("On Arch/Manjaro - try running 'sudo pacman -S bash ffmpeg youtube-dl'")
			Print("On Debian/Ubuntu - try running 'sudo apt install bash ffmpeg youtube-dl'")
			return err
		}
		Excecutables[key] = path
		Print("Excecutable " + key + " found at " + path)
	}
	return nil
}

func ExcecCmd(command string) (string, error) {
	return ExcecProgram("bash", "-c", command)
}
func ExcecCmdToString(command string) (string, error) {
	return ExcecProgramToString("bash", "-c", command)
}

func ExcecProgram(program string, arg ...string) (string, error) {
	args := strings.Join(arg, " ")
	Print("Excecute " + program + " " + args)

	cmd := exec.Command(program, arg...)
	cmd.Dir = WSRoot
	// configure `Stdout` and `Stderr`
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err := cmd.Run()
	// run command
	if err != nil {
		fmt.Println("Error:", err)
	}

	// out := string(ret)
	return "Done Excecute " + program + " " + args, err
}

func ExcecCmdTask(command string, endTask chan bool) (string, error) {
	return ExcecTask("bash", endTask, "-c", command)
}

func ExcecTask(program string, endTask chan bool, arg ...string) (string, error) {
	args := strings.Join(arg, " ")
	Print("Excecute Task " + program + " " + args)

	cmd := exec.Command(program, arg...)
	cmd.Dir = WSRoot
	// configure `Stdout` and `Stderr`
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err := cmd.Start()
	// run command
	if err != nil {
		fmt.Println("Error:", err)
	}

	Kill := <-endTask

	if Kill {
		PrintError(cmd.Process.Signal(os.Kill))
	} else {
		PrintError(cmd.Process.Signal(os.Interrupt))
	}

	// out := string(ret)
	return "Done Excecute Task " + program + " " + args, err
}

func ExcecProgramToString(program string, arg ...string) (string, error) {
	args := strings.Join(arg, " ")
	Print("Excecute " + program + " " + args)

	cmd := exec.Command(program, arg...)
	cmd.Dir = WSRoot
	// configure `Stdout` and `Stderr`
	cmd.Stderr = os.Stdout
	ret, err := cmd.Output()

	out := string(ret)
	return out, err
}
