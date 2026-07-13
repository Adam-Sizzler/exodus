import re

with open("backend/internal/httpapi/subscription/subscription_service.go", "r") as f:
    lines = f.readlines()

new_lines = []
for line in lines:
    if line.startswith("func generateSubscriptionContent("):
        new_lines.append(line.replace("func generateSubscriptionContent(", "func (s *RenderService) generateSubscriptionContent("))
    elif line.startswith("func buildSubscriptionHeaders("):
        new_lines.append(line.replace("func buildSubscriptionHeaders(", "func (s *RenderService) buildSubscriptionHeaders("))
    elif line.startswith("func filterHostsForResponseType("):
        new_lines.append(line.replace("func filterHostsForResponseType(", "func (s *RenderService) filterHostsForResponseType("))
    elif line.startswith("func shuffleHostsIfNeeded("):
        new_lines.append(line.replace("func shuffleHostsIfNeeded(", "func (s *RenderService) shuffleHostsIfNeeded("))
    elif line.startswith("func resolveSubscriptionURL("):
        new_lines.append(line.replace("func resolveSubscriptionURL(", "func (s *RenderService) resolveSubscriptionURL("))
    elif line.startswith("func buildSubscriptionInfoResponse("):
        new_lines.append(line.replace("func buildSubscriptionInfoResponse(", "func (s *RenderService) buildSubscriptionInfoResponse("))
    else:
        new_lines.append(line)

header = """
type RenderService struct {
manager *dbmanager.DatabaseManager
cfg     *config.BackendConfig
}

func NewRenderService(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) *RenderService {
return &RenderService{manager: manager, cfg: cfg}
}

"""

import_end = 0
for i, line in enumerate(new_lines):
    if line.strip() == ")":
        import_end = i
        break

new_lines.insert(import_end + 1, header)

with open("backend/internal/httpapi/subscription/subscription_service.go", "w") as f:
    f.writelines(new_lines)

