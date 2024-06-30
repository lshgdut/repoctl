package util

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	// utilexec "k8s.io/utils/exec"
)

const (
	// ApplyAnnotationsFlag = "save-config"
	DefaultErrorExitCode = 1
	DefaultChunkSize     = 500
)

type debugError interface {
	DebugError() (msg string, args []interface{})
}

// AddSourceToErr adds handleResourcePrefix and source string to error message.
// verb is the string like "creating", "deleting" etc.
// source is the filename or URL to the template file(*.json or *.yaml), or stdin to use to handle the resource.
// func AddSourceToErr(verb string, source string, err error) error {
// 	if source != "" {
// 		if statusError, ok := err.(apierrors.APIStatus); ok {
// 			status := statusError.Status()
// 			status.Message = fmt.Sprintf("error when %s %q: %v", verb, source, status.Message)
// 			return &apierrors.StatusError{ErrStatus: status}
// 		}
// 		return fmt.Errorf("error when %s %q: %v", verb, source, err)
// 	}
// 	return err
// }

var fatalErrHandler = fatal

// BehaviorOnFatal allows you to override the default behavior when a fatal
// error occurs, which is to call os.Exit(code). You can pass 'panic' as a function
// here if you prefer the panic() over os.Exit(1).
func BehaviorOnFatal(f func(string, int)) {
	fatalErrHandler = f
}

// DefaultBehaviorOnFatal allows you to undo any previous override.  Useful in
// tests.
func DefaultBehaviorOnFatal() {
	fatalErrHandler = fatal
}

// fatal prints the message (if provided) and then exits. If V(99) or greater,
// klog.Fatal is invoked for extended information. This is intended for maintainer
// debugging and out of a reasonable range for users.
func fatal(msg string, code int) {
	// nolint:logcheck // Not using the result of klog.V(99) inside the if
	// branch is okay, we just use it to determine how to terminate.
	if klog.V(99).Enabled() {
		klog.FatalDepth(2, msg)
	}
	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
	}
	os.Exit(code)
}

// ErrExit may be passed to CheckError to instruct it to output nothing but exit with
// status code 1.
var ErrExit = fmt.Errorf("exit")

// CheckErr prints a user friendly error to STDERR and exits with a non-zero
// exit code. Unrecognized errors will be printed with an "error: " prefix.
//
// This method is generic to the command in use and may be used by non-Kubectl
// commands.
func CheckErr(err error) {
	checkErr(err, fatalErrHandler)
}

// CheckDiffErr prints a user friendly error to STDERR and exits with a
// non-zero and non-one exit code. Unrecognized errors will be printed
// with an "error: " prefix.
//
// This method is meant specifically for `kubectl diff` and may be used
// by other commands.
func CheckDiffErr(err error) {
	checkErr(err, func(msg string, code int) {
		fatalErrHandler(msg, code+1)
	})
}

// isInvalidReasonStatusError returns true if this is an API Status error with reason=Invalid.
// This is distinct from generic 422 errors we want to fall back to generic error handling.
// func isInvalidReasonStatusError(err error) bool {
// 	if !apierrors.IsInvalid(err) {
// 		return false
// 	}
// 	statusError, isStatusError := err.(*apierrors.StatusError)
// 	if !isStatusError {
// 		return false
// 	}
// 	status := statusError.Status()
// 	return status.Reason == metav1.StatusReasonInvalid
// }

// checkErr formats a given error as a string and calls the passed handleErr
// func with that string and an kubectl exit code.
func checkErr(err error, handleErr func(string, int)) {
	// unwrap aggregates of 1
	// if agg, ok := err.(utilerrors.Aggregate); ok && len(agg.Errors()) == 1 {
	// 	err = agg.Errors()[0]
	// }

	if err == nil {
		return
	}

	switch {
	case err == ErrExit:
		handleErr("", DefaultErrorExitCode)
	// case isInvalidReasonStatusError(err):
	// 	status := err.(*apierrors.StatusError).Status()
	// 	details := status.Details
	// 	s := "The request is invalid"
	// 	if details == nil {
	// 		// if we have no other details, include the message from the server if present
	// 		if len(status.Message) > 0 {
	// 			s += ": " + status.Message
	// 		}
	// 		handleErr(s, DefaultErrorExitCode)
	// 		return
	// 	}
	// 	if len(details.Kind) != 0 || len(details.Name) != 0 {
	// 		s = fmt.Sprintf("The %s %q is invalid", details.Kind, details.Name)
	// 	} else if len(status.Message) > 0 && len(details.Causes) == 0 {
	// 		// only append the message if we have no kind/name details and no causes,
	// 		// since default invalid error constructors duplicate that information in the message
	// 		s += ": " + status.Message
	// 	}

	// 	if len(details.Causes) > 0 {
	// 		errs := statusCausesToAggrError(details.Causes)
	// 		handleErr(MultilineError(s+": ", errs), DefaultErrorExitCode)
	// 	} else {
	// 		handleErr(s, DefaultErrorExitCode)
	// 	}
	// case clientcmd.IsConfigurationInvalid(err):
	// 	handleErr(MultilineError("Error in configuration: ", err), DefaultErrorExitCode)
	default:
		switch err := err.(type) {
		// case *meta.NoResourceMatchError:
		// 	switch {
		// 	case len(err.PartialResource.Group) > 0 && len(err.PartialResource.Version) > 0:
		// 		handleErr(fmt.Sprintf("the server doesn't have a resource type %q in group %q and version %q", err.PartialResource.Resource, err.PartialResource.Group, err.PartialResource.Version), DefaultErrorExitCode)
		// 	case len(err.PartialResource.Group) > 0:
		// 		handleErr(fmt.Sprintf("the server doesn't have a resource type %q in group %q", err.PartialResource.Resource, err.PartialResource.Group), DefaultErrorExitCode)
		// 	case len(err.PartialResource.Version) > 0:
		// 		handleErr(fmt.Sprintf("the server doesn't have a resource type %q in version %q", err.PartialResource.Resource, err.PartialResource.Version), DefaultErrorExitCode)
		// 	default:
		// 		handleErr(fmt.Sprintf("the server doesn't have a resource type %q", err.PartialResource.Resource), DefaultErrorExitCode)
		// 	}
		// case utilerrors.Aggregate:
		// 	handleErr(MultipleErrors(``, err.Errors()), DefaultErrorExitCode)
		// case utilexec.ExitError:
		// 	handleErr(err.Error(), err.ExitStatus())
		default: // for any other error type
			msg, ok := StandardErrorMessage(err)
			if !ok {
				msg = err.Error()
				if !strings.HasPrefix(msg, "error: ") {
					msg = fmt.Sprintf("error: %s", msg)
				}
			}
			handleErr(msg, DefaultErrorExitCode)
		}
	}
}

// func statusCausesToAggrError(scs []metav1.StatusCause) utilerrors.Aggregate {
// 	errs := make([]error, 0, len(scs))
// 	errorMsgs := sets.NewString()
// 	for _, sc := range scs {
// 		// check for duplicate error messages and skip them
// 		msg := fmt.Sprintf("%s: %s", sc.Field, sc.Message)
// 		if errorMsgs.Has(msg) {
// 			continue
// 		}
// 		errorMsgs.Insert(msg)
// 		errs = append(errs, errors.New(msg))
// 	}
// 	return utilerrors.NewAggregate(errs)
// }

// StandardErrorMessage translates common errors into a human readable message, or returns
// false if the error is not one of the recognized types. It may also log extended
// information to klog.
//
// This method is generic to the command in use and may be used by non-Kubectl
// commands.
func StandardErrorMessage(err error) (string, bool) {
	if debugErr, ok := err.(debugError); ok {
		klog.V(4).Infof(debugErr.DebugError())
	}
	// status, isStatus := err.(apierrors.APIStatus)
	// switch {
	// case isStatus:
	// 	switch s := status.Status(); {
	// 	case s.Reason == metav1.StatusReasonUnauthorized:
	// 		return fmt.Sprintf("error: You must be logged in to the server (%s)", s.Message), true
	// 	case len(s.Reason) > 0:
	// 		return fmt.Sprintf("Error from server (%s): %s", s.Reason, err.Error()), true
	// 	default:
	// 		return fmt.Sprintf("Error from server: %s", err.Error()), true
	// 	}
	// case apierrors.IsUnexpectedObjectError(err):
	// 	return fmt.Sprintf("Server returned an unexpected response: %s", err.Error()), true
	// }
	switch t := err.(type) {
	case *url.Error:
		klog.V(4).Infof("Connection error: %s %s: %v", t.Op, t.URL, t.Err)
		switch {
		case strings.Contains(t.Err.Error(), "connection refused"):
			host := t.URL
			if server, err := url.Parse(t.URL); err == nil {
				host = server.Host
			}
			return fmt.Sprintf("The connection to the server %s was refused - did you specify the right host or port?", host), true
		}
		return fmt.Sprintf("Unable to connect to the server: %v", t.Err), true
	}
	return "", false
}

// MultilineError returns a string representing an error that splits sub errors into their own
// lines. The returned string will end with a newline.
func MultilineError(prefix string, err error) string {
	// if agg, ok := err.(utilerrors.Aggregate); ok {
	// 	errs := utilerrors.Flatten(agg).Errors()
	// 	buf := &bytes.Buffer{}
	// 	switch len(errs) {
	// 	case 0:
	// 		return fmt.Sprintf("%s%v\n", prefix, err)
	// 	case 1:
	// 		return fmt.Sprintf("%s%v\n", prefix, messageForError(errs[0]))
	// 	default:
	// 		fmt.Fprintln(buf, prefix)
	// 		for _, err := range errs {
	// 			fmt.Fprintf(buf, "* %v\n", messageForError(err))
	// 		}
	// 		return buf.String()
	// 	}
	// }
	return fmt.Sprintf("%s%s\n", prefix, err)
}

// PrintErrorWithCauses prints an error's kind, name, and each of the error's causes in a new line.
// The returned string will end with a newline.
// Returns true if a case exists to handle the error type, or false otherwise.
func PrintErrorWithCauses(err error, errOut io.Writer) bool {
	// switch t := err.(type) {
	// case *apierrors.StatusError:
	// 	errorDetails := t.Status().Details
	// 	if errorDetails != nil {
	// 		fmt.Fprintf(errOut, "error: %s %q is invalid\n\n", errorDetails.Kind, errorDetails.Name)
	// 		for _, cause := range errorDetails.Causes {
	// 			fmt.Fprintf(errOut, "* %s: %s\n", cause.Field, cause.Message)
	// 		}
	// 		return true
	// 	}
	// }

	fmt.Fprintf(errOut, "error: %v\n", err)
	return false
}

// MultipleErrors returns a newline delimited string containing
// the prefix and referenced errors in standard form.
func MultipleErrors(prefix string, errs []error) string {
	buf := &bytes.Buffer{}
	for _, err := range errs {
		fmt.Fprintf(buf, "%s%v\n", prefix, messageForError(err))
	}
	return buf.String()
}

// messageForError returns the string representing the error.
func messageForError(err error) string {
	msg, ok := StandardErrorMessage(err)
	if !ok {
		msg = err.Error()
	}
	return msg
}

func UsageErrorf(cmd *cobra.Command, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s\nSee '%s -h' for help and examples", msg, cmd.CommandPath())
}

func IsFilenameSliceEmpty(filenames []string, directory string) bool {
	return len(filenames) == 0 && directory == ""
}

func GetFlagString(cmd *cobra.Command, flag string) string {
	s, err := cmd.Flags().GetString(flag)
	if err != nil {
		klog.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return s
}

// GetFlagStringSlice can be used to accept multiple argument with flag repetition (e.g. -f arg1,arg2 -f arg3 ...)
func GetFlagStringSlice(cmd *cobra.Command, flag string) []string {
	s, err := cmd.Flags().GetStringSlice(flag)
	if err != nil {
		klog.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return s
}

// GetFlagStringArray can be used to accept multiple argument with flag repetition (e.g. -f arg1 -f arg2 ...)
func GetFlagStringArray(cmd *cobra.Command, flag string) []string {
	s, err := cmd.Flags().GetStringArray(flag)
	if err != nil {
		klog.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return s
}

func GetFlagBool(cmd *cobra.Command, flag string) bool {
	b, err := cmd.Flags().GetBool(flag)
	if err != nil {
		klog.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return b
}

// Assumes the flag has a default value.
func GetFlagInt(cmd *cobra.Command, flag string) int {
	i, err := cmd.Flags().GetInt(flag)
	if err != nil {
		klog.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return i
}

// Assumes the flag has a default value.
func GetFlagInt32(cmd *cobra.Command, flag string) int32 {
	i, err := cmd.Flags().GetInt32(flag)
	if err != nil {
		klog.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return i
}

// Assumes the flag has a default value.
func GetFlagInt64(cmd *cobra.Command, flag string) int64 {
	i, err := cmd.Flags().GetInt64(flag)
	if err != nil {
		klog.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return i
}

func GetFlagDuration(cmd *cobra.Command, flag string) time.Duration {
	d, err := cmd.Flags().GetDuration(flag)
	if err != nil {
		klog.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return d
}

// AddDryRunFlag adds dry-run flag to a command. Usually used by mutations.
func AddDryRunFlag(cmd *cobra.Command) {
	cmd.Flags().Bool(
		"dry-run",
		false,
		"If true, only print the object that would be sent, without sending it. ",
	)
	cmd.Flags().Lookup("dry-run").NoOptDefVal = "true"
}
