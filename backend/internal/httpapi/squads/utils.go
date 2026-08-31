package squads

import (
	"database/sql"

	"exodus/internal/db"
	"exodus/internal/httpapi/shared"
)

func scanInternalSquad(scanner shared.RowScanner) (InternalSquad, error) {
	var squad InternalSquad
	var viewPosition sql.NullInt64
	var tags db.StringArray

	err := scanner.Scan(
		&squad.UUID,
		&viewPosition,
		&squad.Name,
		&tags,
		&squad.CreatedAt,
		&squad.UpdatedAt,
	)
	if err != nil {
		return squad, err
	}

	if viewPosition.Valid {
		squad.ViewPosition = int(viewPosition.Int64)
	}

	squad.Tags = tags.Slice()
	if squad.Tags == nil {
		squad.Tags = []string{}
	}

	return squad, nil
}

func buildInternalSquadResponse(squad InternalSquad, membersCount int, inbounds []InternalSquadInboundAPI) InternalSquadAPI {
	if inbounds == nil {
		inbounds = []InternalSquadInboundAPI{}
	}
	tags := squad.Tags
	if tags == nil {
		tags = []string{}
	}
	return InternalSquadAPI{
		UUID:         squad.UUID,
		ViewPosition: squad.ViewPosition,
		Name:         squad.Name,
		Tags:         tags,
		Info: InternalSquadInfo{
			MembersCount:  membersCount,
			InboundsCount: len(inbounds),
		},
		Inbounds:  inbounds,
		CreatedAt: squad.CreatedAt,
		UpdatedAt: squad.UpdatedAt,
	}
}
