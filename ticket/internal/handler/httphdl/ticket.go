package httphdl

import (
	"net/http"
	"ticket/internal/core/domain"
	"ticket/pkg/errors"

	"github.com/labstack/echo/v4"
)

func (hdl *HTTPHandler) CreateTicket(c echo.Context) error {
	ctx := c.Request().Context()
	auth, ok := ctx.Value(domain.ContextKeyAuth).(*domain.AuthenticationResponse)
	if !ok {
		return errors.ErrUnauthorizedContext
	}

	var request domain.TicketRequest
	if err := c.Bind(&request); err != nil {
		return errors.NewResponseErr(err)
	}

	if err := hdl.validator.StrcutWithTranslateError(request); err != nil {
		return errors.NewResponseErr(errors.ErrValidation, errors.OptionErr(err))
	}

	resp, err := hdl.ticketService.CreateTicket(ctx, &auth.User, &request)
	if err != nil {
		return errors.NewResponseErr(err)
	}
	return c.JSON(http.StatusCreated, resp)
}

func (hdl *HTTPHandler) GetTicketList(c echo.Context) error {
	ctx := c.Request().Context()
	var query domain.GetTicketListQuery
	if err := c.Bind(&query); err != nil {
		return errors.NewResponseErr(err)
	}

	if err := hdl.validator.StrcutWithTranslateError(query); err != nil {
		return errors.NewResponseErr(errors.ErrValidation, errors.OptionErr(err))
	}

	resp, err := hdl.ticketService.GetTicketList(ctx, query)
	if err != nil {
		return errors.NewResponseErr(err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (hdl *HTTPHandler) GetTicket(c echo.Context) error {
	ctx := c.Request().Context()
	var query domain.GetTicketQuery
	if err := c.Bind(&query); err != nil {
		return errors.NewResponseErr(err)
	}

	if err := hdl.validator.StrcutWithTranslateError(query); err != nil {
		return errors.NewResponseErr(errors.ErrValidation, errors.OptionErr(err))
	}

	resp, err := hdl.ticketService.GetTicketByID(ctx, query.ID)
	if err != nil {
		return errors.NewResponseErr(err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (hdl *HTTPHandler) UpdateTicket(c echo.Context) error {
	ctx := c.Request().Context()
	auth, ok := ctx.Value(domain.ContextKeyAuth).(*domain.AuthenticationResponse)
	if !ok {
		return errors.ErrUnauthorizedContext
	}

	var request domain.UpdateTicketRequest
	if err := c.Bind(&request); err != nil {
		return errors.NewResponseErr(err)
	}

	if err := hdl.validator.StrcutWithTranslateError(request); err != nil {
		return errors.NewResponseErr(errors.ErrValidation, errors.OptionErr(err))
	}

	resp, err := hdl.ticketService.UpdateTicket(ctx, &auth.User, &request)
	if err != nil {
		return errors.NewResponseErr(err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (hdl *HTTPHandler) DeleteTicket(c echo.Context) error {
	ctx := c.Request().Context()
	var query domain.DeleteTicketQuery
	if err := c.Bind(&query); err != nil {
		return errors.NewResponseErr(err)
	}

	if err := hdl.validator.StrcutWithTranslateError(query); err != nil {
		return errors.NewResponseErr(errors.ErrValidation, errors.OptionErr(err))
	}

	resp, err := hdl.ticketService.DeleteTicket(ctx, query.ID)
	if err != nil {
		return errors.NewResponseErr(err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (hdl *HTTPHandler) GetTicketHistoryList(c echo.Context) error {
	ctx := c.Request().Context()
	var query domain.GetTicketHistoryListQuery
	if err := c.Bind(&query); err != nil {
		return errors.NewResponseErr(err)
	}

	if err := hdl.validator.StrcutWithTranslateError(query); err != nil {
		return errors.NewResponseErr(errors.ErrValidation, errors.OptionErr(err))
	}

	resp, err := hdl.ticketService.GetTicketHistoryList(ctx, query)
	if err != nil {
		return errors.NewResponseErr(err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (hdl *HTTPHandler) CreateTicketComment(c echo.Context) error {
	ctx := c.Request().Context()
	auth, ok := ctx.Value(domain.ContextKeyAuth).(*domain.AuthenticationResponse)
	if !ok {
		return errors.ErrUnauthorizedContext
	}

	var request domain.TicketCommentRequest
	if err := c.Bind(&request); err != nil {
		return errors.NewResponseErr(err)
	}

	if err := hdl.validator.StrcutWithTranslateError(request); err != nil {
		return errors.NewResponseErr(errors.ErrValidation, errors.OptionErr(err))
	}

	resp, err := hdl.ticketService.CreateTicketComment(ctx, &auth.User, &request)
	if err != nil {
		return errors.NewResponseErr(err)
	}
	return c.JSON(http.StatusCreated, resp)
}

func (hdl *HTTPHandler) GetTicketCommentList(c echo.Context) error {
	ctx := c.Request().Context()
	var query domain.GetTicketCommentListQuery
	if err := c.Bind(&query); err != nil {
		return errors.NewResponseErr(err)
	}

	if err := hdl.validator.StrcutWithTranslateError(query); err != nil {
		return errors.NewResponseErr(errors.ErrValidation, errors.OptionErr(err))
	}

	resp, err := hdl.ticketService.GetTicketCommentList(ctx, query)
	if err != nil {
		return errors.NewResponseErr(err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (hdl *HTTPHandler) UpdateTicketComment(c echo.Context) error {
	ctx := c.Request().Context()
	auth, ok := ctx.Value(domain.ContextKeyAuth).(*domain.AuthenticationResponse)
	if !ok {
		return errors.ErrUnauthorizedContext
	}

	var request domain.UpdateTicketCommentRequest
	if err := c.Bind(&request); err != nil {
		return errors.NewResponseErr(err)
	}

	if err := hdl.validator.StrcutWithTranslateError(request); err != nil {
		return errors.NewResponseErr(errors.ErrValidation, errors.OptionErr(err))
	}

	resp, err := hdl.ticketService.UpdateTicketComment(ctx, &auth.User, &request)
	if err != nil {
		return errors.NewResponseErr(err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (hdl *HTTPHandler) DeleteTicketComment(c echo.Context) error {
	ctx := c.Request().Context()
	auth, ok := ctx.Value(domain.ContextKeyAuth).(*domain.AuthenticationResponse)
	if !ok {
		return errors.ErrUnauthorizedContext
	}

	var query domain.DeleteTicketCommentQuery
	if err := c.Bind(&query); err != nil {
		return errors.NewResponseErr(err)
	}

	if err := hdl.validator.StrcutWithTranslateError(query); err != nil {
		return errors.NewResponseErr(errors.ErrValidation, errors.OptionErr(err))
	}

	resp, err := hdl.ticketService.DeleteTicketComment(ctx, &auth.User, query.ID)
	if err != nil {
		return errors.NewResponseErr(err)
	}
	return c.JSON(http.StatusOK, resp)
}
