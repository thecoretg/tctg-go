package psa

import (
	"context"
	"fmt"
)

const (
	ticketLinkBase = "https://na.myconnectwise.net/v4_6_release/services/system_io/Service/fv_sr100_request.rails?service_recid="
)

func ticketIdEndpoint(ticketId int) string {
	return fmt.Sprintf("service/tickets/%d", ticketId)
}

func notesEndpoint(ticketId int) string {
	return fmt.Sprintf("%s/notes", ticketIdEndpoint(ticketId))
}

func allNotesEndpoint(ticketId int) string {
	return fmt.Sprintf("%s/allNotes", ticketIdEndpoint(ticketId))
}

func specificNoteEndpoint(ticketId, noteId int) string {
	return fmt.Sprintf("%s/notes/%d", ticketIdEndpoint(ticketId), noteId)
}

func (c *Client) PostTicket(ctx context.Context, ticket *Ticket) (*Ticket, error) {
	return Post[Ticket](ctx, c, "service/tickets", ticket)
}

func (c *Client) ListTickets(ctx context.Context, params map[string]string) ([]Ticket, error) {
	return GetMany[Ticket](ctx, c, "service/tickets", params)
}

func (c *Client) GetTicket(ctx context.Context, ticketID int, params map[string]string) (*Ticket, error) {
	return GetOne[Ticket](ctx, c, ticketIdEndpoint(ticketID), params)
}

func (c *Client) PutTicket(ctx context.Context, ticketID int, ticket *Ticket) (*Ticket, error) {
	return Put[Ticket](ctx, c, ticketIdEndpoint(ticketID), ticket)
}

func (c *Client) PatchTicket(ctx context.Context, ticketID int, patchOps []PatchOp) (*Ticket, error) {
	return Patch[Ticket](ctx, c, ticketIdEndpoint(ticketID), patchOps)
}

func (c *Client) DeleteTicket(ctx context.Context, ticketID int) error {
	return Delete(ctx, c, ticketIdEndpoint(ticketID))
}

// ListServiceTicketNotesAll gets all ticket notes, regardless of if they have a time entry.
//
// This is most likely the one you want to use unless you consistently uncheck the time entry box.
func (c *Client) ListServiceTicketNotesAll(ctx context.Context, params map[string]string, ticketID int) ([]ServiceTicketNoteAll, error) {
	return GetMany[ServiceTicketNoteAll](ctx, c, allNotesEndpoint(ticketID), params)
}

func (c *Client) PostServiceTicketNote(ctx context.Context, ticketNote *ServiceTicketNote, ticketID int) (*ServiceTicketNote, error) {
	return Post[ServiceTicketNote](ctx, c, notesEndpoint(ticketID), ticketNote)
}

// ListServiceTicketNotes gets all notes that are not time entry.
//
// Not recommended since you will probably get what you need through ListServiceTicketNotesAll.
func (c *Client) ListServiceTicketNotes(ctx context.Context, params map[string]string, ticketID int) ([]ServiceTicketNote, error) {
	return GetMany[ServiceTicketNote](ctx, c, notesEndpoint(ticketID), params)
}

func (c *Client) GetServiceTicketNote(ctx context.Context, noteID int, params map[string]string, ticketID int) (*ServiceTicketNote, error) {
	return GetOne[ServiceTicketNote](ctx, c, specificNoteEndpoint(ticketID, noteID), params)
}

func (c *Client) PutServiceTicketNote(ctx context.Context, noteID int, ticketNote *ServiceTicketNote, ticketID int) (*ServiceTicketNote, error) {
	return Put[ServiceTicketNote](ctx, c, specificNoteEndpoint(ticketID, noteID), ticketNote)
}

func (c *Client) PatchServiceTicketNote(ctx context.Context, noteID int, patchOps []PatchOp, ticketID int) (*ServiceTicketNote, error) {
	return Patch[ServiceTicketNote](ctx, c, specificNoteEndpoint(ticketID, noteID), patchOps)
}

func (c *Client) DeleteServiceTicketNote(ctx context.Context, noteID int, ticketID int) error {
	return Delete(ctx, c, specificNoteEndpoint(ticketID, noteID))
}

func (c *Client) GetMostRecentTicketNote(ctx context.Context, ticketID int) (*ServiceTicketNote, error) {
	p := map[string]string{
		"orderBy":  "id desc",
		"pageSize": "1000",
	}

	notes, err := c.ListServiceTicketNotesAll(ctx, p, ticketID)
	if err != nil {
		return nil, fmt.Errorf("listing service notes: %w", err)
	}

	if len(notes) == 0 {
		return nil, nil
	}

	note, err := c.GetServiceTicketNote(ctx, notes[0].ID, nil, ticketID)
	if err != nil {
		return nil, fmt.Errorf("getting details for note: %w", err)
	}

	return note, nil
}

func MarkdownInternalTicketLink(ticketID int, companyID string) string {
	return fmt.Sprintf("[%d](%s)", ticketID, InternalTicketLink(ticketID, companyID))
}

func InternalTicketLink(ticketID int, companyID string) string {
	return fmt.Sprintf("%s%d&companyName=%s", ticketLinkBase, ticketID, companyID)
}
