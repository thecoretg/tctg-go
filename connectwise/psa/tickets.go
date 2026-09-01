package psa

import (
	"context"
	"fmt"
)

const (
	ticketLinkBase = "https://na.myconnectwise.net/v4_6_release/services/system_io/Service/fv_sr100_request.rails?service_recid="
)

func ticketIDEndpoint(ticketID int) string {
	return fmt.Sprintf("service/tickets/%d", ticketID)
}

func notesEndpoint(ticketID int) string {
	return fmt.Sprintf("%s/notes", ticketIDEndpoint(ticketID))
}

func allNotesEndpoint(ticketID int) string {
	return fmt.Sprintf("%s/allNotes", ticketIDEndpoint(ticketID))
}

func specificNoteEndpoint(ticketID, noteID int) string {
	return fmt.Sprintf("%s/notes/%d", ticketIDEndpoint(ticketID), noteID)
}

func (c *Client) PostTicket(ctx context.Context, ticket *Ticket) (*Ticket, error) {
	return c.Post[Ticket](ctx, "service/tickets", ticket)
}

func (c *Client) ListTickets(ctx context.Context, params map[string]string, opts ...ListOption) ([]Ticket, error) {
	return c.GetMany[Ticket](ctx, "service/tickets", params, opts...)
}

func (c *Client) GetTicket(ctx context.Context, ticketID int, params map[string]string) (*Ticket, error) {
	return c.Get[Ticket](ctx, ticketIDEndpoint(ticketID), params)
}

func (c *Client) PutTicket(ctx context.Context, ticketID int, ticket *Ticket) (*Ticket, error) {
	return c.Put[Ticket](ctx, ticketIDEndpoint(ticketID), ticket)
}

func (c *Client) PatchTicket(ctx context.Context, ticketID int, patchOps []PatchOp) (*Ticket, error) {
	return c.Patch[Ticket](ctx, ticketIDEndpoint(ticketID), patchOps)
}

func (c *Client) DeleteTicket(ctx context.Context, ticketID int) error {
	return c.Delete(ctx, ticketIDEndpoint(ticketID))
}

// ListServiceTicketNotesAll gets all ticket notes, regardless of if they have a time entry.
//
// This is most likely the one you want to use unless you consistently uncheck the time entry box.
func (c *Client) ListServiceTicketNotesAll(ctx context.Context, params map[string]string, ticketID int, opts ...ListOption) ([]ServiceTicketNoteAll, error) {
	return c.GetMany[ServiceTicketNoteAll](ctx, allNotesEndpoint(ticketID), params, opts...)
}

func (c *Client) PostServiceTicketNote(ctx context.Context, ticketNote *ServiceTicketNote, ticketID int) (*ServiceTicketNote, error) {
	return c.Post[ServiceTicketNote](ctx, notesEndpoint(ticketID), ticketNote)
}

// ListServiceTicketNotes gets all notes that are not time entry.
//
// Not recommended since you will probably get what you need through ListServiceTicketNotesAll.
func (c *Client) ListServiceTicketNotes(ctx context.Context, params map[string]string, ticketID int, opts ...ListOption) ([]ServiceTicketNote, error) {
	return c.GetMany[ServiceTicketNote](ctx, notesEndpoint(ticketID), params, opts...)
}

func (c *Client) GetServiceTicketNote(ctx context.Context, noteID int, params map[string]string, ticketID int) (*ServiceTicketNote, error) {
	return c.Get[ServiceTicketNote](ctx, specificNoteEndpoint(ticketID, noteID), params)
}

func (c *Client) PutServiceTicketNote(ctx context.Context, noteID int, ticketNote *ServiceTicketNote, ticketID int) (*ServiceTicketNote, error) {
	return c.Put[ServiceTicketNote](ctx, specificNoteEndpoint(ticketID, noteID), ticketNote)
}

func (c *Client) PatchServiceTicketNote(ctx context.Context, noteID int, patchOps []PatchOp, ticketID int) (*ServiceTicketNote, error) {
	return c.Patch[ServiceTicketNote](ctx, specificNoteEndpoint(ticketID, noteID), patchOps)
}

func (c *Client) DeleteServiceTicketNote(ctx context.Context, noteID int, ticketID int) error {
	return c.Delete(ctx, specificNoteEndpoint(ticketID, noteID))
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
