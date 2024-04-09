package maintainer

import (
	"context"
	"fmt"
	"strings"

	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/sender"
)

const AddAdminCommand = "/addAdmin"
const RemoveAdminCommand = "/removeAdmin"

func NewAddAdmin(deps *fsm.Deps) *AddAdmin {
	return &AddAdmin{
		Admin: Admin{deps: deps},
	}
}

type AddAdmin struct {
	Admin
}

func (s *AddAdmin) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	successAdminText := "Successful give permissions to user %s"
	successUserText := "Maintainer give you manager permissions"
	return s.changeManagerPermissions(ctx, state, AddAdminCommand, true, successAdminText, successUserText)
}

func NewRemoveAdmin(deps *fsm.Deps) *RemoveAdmin {
	return &RemoveAdmin{
		Admin: Admin{deps: deps},
	}
}

type RemoveAdmin struct {
	Admin
}

func (s *RemoveAdmin) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	successAdminText := "Successful remove permissions to user %s"
	successUserText := "Maintainer remove your manager permissions"
	return s.changeManagerPermissions(ctx, state, RemoveAdminCommand, false, successAdminText, successUserText)
}

type Admin struct {
	deps *fsm.Deps
}

func (s *Admin) changeManagerPermissions(
	ctx context.Context,
	state fsm.State,
	command string,
	isManager bool,
	successAdminText string,
	successUserText string,
) (context.Context, fsm.Machine, fsm.State, error) {
	if state.Update.Message == nil {
		return ctx, nil, state, fmt.Errorf("unexpected part. ")
	}
	if s.deps.Cfg.MaintainerChatID == state.User.ID && !state.User.IsMaintainer {
		smc := sender.NewSenderMachine(s.deps, state.Update.Message.Chat.ID, "You dont have enough permissions for this action.", 0)
		return ctx, smc, state, nil
	}

	targetUsername := strings.TrimPrefix(state.Update.Message.Text, command)
	targetUsername = strings.TrimSpace(targetUsername)
	targetUsername = strings.TrimPrefix(targetUsername, "@")

	targetUser, err := s.deps.Repo.GetUserByUsername(ctx, targetUsername)
	if err != nil {
		return ctx, nil, state, fmt.Errorf("repo.GetUserByUsername() err: %w", err)
	}
	if targetUser.ID == 0 {
		smc := sender.NewSenderMachine(s.deps, state.Update.Message.Chat.ID, fmt.Sprintf("I don't know the user %s", targetUsername), 0)
		return ctx, smc, state, nil
	}
	err = s.deps.Repo.SetUserIsManager(ctx, targetUser.ID, isManager)
	if err != nil {
		return ctx, nil, state, fmt.Errorf("repo.SaveUser() err: %w", err)
	}

	combineMachine := fsm.NewCombine(nil,
		sender.NewSenderMachine(s.deps, state.Update.Message.Chat.ID, fmt.Sprintf(successAdminText, targetUsername), 0),
		sender.NewSenderMachine(s.deps, targetUser.ID, successUserText, 0),
	)

	return ctx, combineMachine, state, nil
}
