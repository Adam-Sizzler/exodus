package squads

import (
	"database/sql"

	"exodus/internal/httpapi/shared"
)

func scanInternalSquad(scanner shared.RowScanner) (InternalSquad, error) {
	var squad InternalSquad
	var viewPosition sql.NullInt64

	err := scanner.Scan(
		&squad.UUID,
		&viewPosition,
		&squad.Name,
		&squad.CreatedAt,
		&squad.UpdatedAt,
	)
	if err != nil {
		return squad, err
	}

	if viewPosition.Valid {
		squad.ViewPosition = int(viewPosition.Int64)
	}

	return squad, nil
}

func buildInternalSquadResponse(squad InternalSquad, membersCount int, inbounds []InternalSquadInboundAPI) InternalSquadAPI {
	if inbounds == nil {
		inbounds = []InternalSquadInboundAPI{}
	}
	return InternalSquadAPI{
		UUID:         squad.UUID,
		ViewPosition: squad.ViewPosition,
		Name:         squad.Name,
		Info: InternalSquadInfo{
			MembersCount:  membersCount,
			InboundsCount: len(inbounds),
		},
		Inbounds:  inbounds,
		CreatedAt: squad.CreatedAt,
		UpdatedAt: squad.UpdatedAt,
	}
}
