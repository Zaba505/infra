package main

import "dagger/rebuild/run"

// phases is the rebuild's phase order, and the one place it is fixed. up
// applies and verifies the phases in this order, and down tears them down in
// reverse, so the first phase, phase 0, is torn down last.
//
// A phase joins the order when its proposal's stories land it:
//
//	0   the state backend                              Proposal: definitions, values and secrets
//	1a  the cloud anchor and the edge                  Proposal: cloud anchor and edge
//	1b  the operations tunnel                          Proposal: operations tunnel
//	1c  home-lab hosts                                 Proposal: home-lab hosts
//	2   tenant compute, then the offerings tenants use Proposal: tenant compute, then each wave-5 proposal
//	3   backup, and observability                      Proposal: backup, restore and export; Proposal: observability
//	4   the canary tenant                              Proposal: canary tenant
//
// No phase has landed yet, so the order is empty. A run against any target
// then passes the preflight check as a rebuild from nothing and finishes with
// nothing to do.
func (m *Rebuild) phases() []run.Phase {
	return nil
}
