// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package auth

import (
	"github.com/playbymail/ottoapp/backend/domains"
)

// policy helpers return true if the actor is permitted to take an action.

func (s *Service) CanAuthenticate(actor *domains.Actor) bool {
	// Only "user" role principals may authenticate via credentials.
	// Sysop and service actors are always blocked from password-based login.
	return actor.IsUser() && !(actor.IsSysop() || actor.IsService())
}

// CanCreateDocuments returns true if the actor can create documents.
func (s *Service) CanCreateDocuments(actor *domains.Actor) bool {
	// todo: implement a real policy
	return actor.IsValid()
}

// CanCreateTarget checks if actor can create new users.
func (s *Service) CanCreateTarget(actor *domains.Actor) bool {
	if actor.IsSysop() {
		// sysop can create users
		return true
	}
	// from here on, sysop is impossible

	// only admins can create users
	if !actor.IsAdmin() {
		return false
	}

	return true
}

// CanEditTarget checks if actor can edit target user's profile.
// Rules: user can edit self, admin can edit non-admins (excluding sysop).
func (s *Service) CanEditTarget(actor, target *domains.Actor) bool {
	if target.IsSysop() {
		// no one can edit sysop, not even syspo
		return false
	} else if actor.IsSysop() {
		// sysop can edit anyone
		return true
	}
	// from here on, sysop is impossible

	// user can edit themselves
	if actor.ID == target.ID {
		return true
	}

	// only admins can edit other users
	if !actor.IsAdmin() {
		return false
	}

	// admin cannot edit other admins
	if target.IsAdmin() {
		return false
	}

	return true
}

// CanEditTargetUsername checks if actor can edit target's username.
func (s *Service) CanEditTargetUsername(actor, target *domains.Actor) bool {
	if target.IsSysop() {
		// no one can edit sysop, not even sysop
		return false
	} else if actor.IsSysop() {
		// sysop can edit anyone
		return true
	}
	// from here on, sysop is impossible

	// only admins can edit usernames
	if !actor.IsAdmin() {
		return false
	}

	// admin cannot edit other admins
	if target.IsAdmin() {
		return false
	}

	return true
}

func (s *Service) CanListUsers(actor *domains.Actor) bool {
	if actor.IsSysop() {
		// sysop can edit anyone
		return true
	}
	// from here on, sysop is impossible

	// only admins can list users
	if !actor.IsAdmin() {
		return false
	}

	return true
}

// CanManageTargetRoles checks if actor can modify target's roles.
func (s *Service) CanManageTargetRoles(actor, target *domains.Actor) bool {
	if target.IsSysop() {
		// no one can manage sysop's roles, not even sysop
		return false
	} else if actor.IsSysop() {
		// sysop can manage anyone's roles
		return true
	}
	// from here on, sysop is impossible

	if target.IsService() {
		// no one can manage service roles, not even services
		return false
	}

	// only admins can manage roles
	if !actor.IsAdmin() {
		return false
	}

	// admin can't manage roles for another admin (even themselves)
	if target.IsAdmin() {
		return false
	}

	return true
}

// CanResetTargetCredentials checks if actor can reset target's credentials.
// Only admins can reset passwords for non-admins.
func (s *Service) CanResetTargetCredentials(actor, target *domains.Actor) bool {
	if target.IsSysop() {
		// no one can reset sysop's credentials, not even sysop
		// (this is to prevent sysop from ever being able to log in)
		return false
	} else if actor.IsSysop() {
		// sysop can reset anyone's credentials
		return true
	}
	// from here on, sysop is impossible

	// users can reset their own credentials
	if actor.ID == target.ID {
		return true
	}

	// only admins can reset other users' credentials
	if !actor.IsAdmin() {
		return false
	}

	// admin can't reset other admin's credentials
	if target.IsAdmin() {
		return false
	}

	return true
}

func (s *Service) CanShutdownServer(actor *domains.Actor) bool {
	if actor.IsSysop() {
		// sysop can server
		return true
	}
	// from here on, sysop is impossible

	if actor.IsService() {
		// services can shut down server
		return true
	}

	return false
}

func (s *Service) CanUpdateTargetCredentials(actor, target *domains.Actor) bool {
	if target.IsSysop() {
		// no one is allowed to change the credentials for sysop
		// (this is to prevent sysop from ever being able to authenticate)
		return false
	} else if actor.IsSysop() {
		// sysop can change everyone else's credentials
		return true
	}
	// from here down, sysop is impossible

	if actor.ID == target.ID {
		// everyone else can change their own credentials
		return true
	}
	// all other actor/target combinations require admin role
	if !actor.IsAdmin() {
		return false
	}
	// last check is for admins - they can't update the credentials
	// for another admin. this is meant to prevent a rogue admin
	// from blocking access to the rest of the team.
	if target.IsAdmin() {
		return false
	}

	return true
}

// CanViewTarget returns true if actor can view target's profile.
// Rules: user can view self, admin can view all (excluding sysop)
func (s *Service) CanViewTarget(actor, target *domains.Actor) bool {
	if target.IsSysop() {
		// no one can view sysop, not even sysop
		return false
	} else if actor.IsSysop() {
		// sysop can view anyone
		return true
	}
	// from here on, sysop is impossible

	// users can view themselves
	if actor.ID == target.ID {
		return true
	}

	// only admins can view others
	if !actor.IsAdmin() {
		return false
	}

	return true
}
