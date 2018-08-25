package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/almostmoore/gbquestion/question"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var client question.QuestionsClient
var upsertCmd, listCmd, deleteCmd, viewCmd *cobra.Command

func initClient(cmd *cobra.Command, args []string) error {
	conn, err := grpc.Dial(os.Getenv("LISTEN"), grpc.WithInsecure())
	if err != nil {
		return err
	}

	client = question.NewQuestionsClient(conn)
	return nil
}

func upsert(cmd *cobra.Command, args []string) error {
	q := &question.Question{}
	q.Id, _ = cmd.Flags().GetUint64("id")
	q.Text, _ = cmd.Flags().GetString("text")
	q.IsActive, _ = cmd.Flags().GetBool("active")
	q.IsGood, _ = cmd.Flags().GetBool("good")

	q, err := client.Put(context.Background(), q)
	if err != nil {
		return fmt.Errorf("Couldn't send a question: %v", err)
	}

	renderQuestions([]*question.Question{q})

	return nil
}

func list(cmd *cobra.Command, args []string) error {
	filter := &question.Filter{}
	filter.Limit, _ = cmd.Flags().GetInt32("limit")
	filter.IsActive, _ = cmd.Flags().GetBool("active")
	filter.Offset, _ = cmd.Flags().GetInt32("offset")

	l, err := client.List(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("Couldn't fetch a list of questions: %v", err)
	}

	renderQuestions(l.Questions)
	return nil
}

func view(cmd *cobra.Command, args []string) error {
	idRequest := &question.IdRequest{}
	idRequest.Id, _ = cmd.Flags().GetUint64("id")

	q, err := client.Get(context.Background(), idRequest)
	if err != nil {
		return fmt.Errorf("Unable to fetch a question: %v", err)
	}

	renderQuestions([]*question.Question{q})
	return nil
}

func delete(cmd *cobra.Command, args []string) error {
	idRequest := &question.IdRequest{}
	idRequest.Id, _ = cmd.Flags().GetUint64("id")

	_, err := client.Delete(context.Background(), idRequest)
	if err != nil {
		return fmt.Errorf("Unable to delete a question: %v", err)
	}

	fmt.Printf("Question %d was deleted successfully\n", idRequest.Id)
	return nil
}

func renderQuestions(questions []*question.Question) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Text", "Is Active", "Is Good"})

	for _, q := range questions {
		table.Append([]string{
			strconv.FormatUint(q.Id, 10),
			q.Text,
			strconv.FormatBool(q.IsActive),
			strconv.FormatBool(q.IsGood),
		})
	}

	table.Render()
}

func init() {
	upsertCmd = &cobra.Command{
		Use:     "upsert",
		Short:   "Add or update question",
		PreRunE: initClient,
		RunE:    upsert,
	}

	upsertCmd.Flags().StringP("text", "t", "", "Text of the question")
	upsertCmd.Flags().Uint64P("id", "", 0, "ID of the question")
	upsertCmd.Flags().BoolP("active", "a", true, "Flag of activity")
	upsertCmd.Flags().BoolP("good", "g", true, "Is it a good answer?")

	listCmd = &cobra.Command{
		Use:     "list",
		Short:   "Show list of questions",
		PreRunE: initClient,
		RunE:    list,
	}

	listCmd.Flags().Int32P("limit", "l", 100, "Limit of questions")
	listCmd.Flags().Int32P("offset", "o", 0, "Offset from the start")
	listCmd.Flags().BoolP("active", "a", true, "Show only active or disabled question")

	viewCmd = &cobra.Command{
		Use:     "view",
		Short:   "Show one question",
		PreRunE: initClient,
		RunE:    view,
	}

	viewCmd.Flags().Uint64P("id", "", 0, "Id of a question")

	deleteCmd = &cobra.Command{
		Use:     "delete",
		Short:   "Delete one question",
		PreRunE: initClient,
		RunE:    delete,
	}

	deleteCmd.Flags().Uint64P("id", "", 0, "Id of a question")
}
