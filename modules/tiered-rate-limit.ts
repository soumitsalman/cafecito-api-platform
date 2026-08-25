import { ZuploContext, ZuploRequest } from "@zuplo/runtime";

const FREE_REQUESTS_PER_MINUTE = 100;
const BALLER_REQUESTS_PER_MINUTE = 1000;

export function rateLimit(request: ZuploRequest, _context: ZuploContext) {
  const sub = request.user?.sub;
  const plan = (request.user?.data as Record<string, unknown> | undefined)
    ?.subscription_plan;

  if (plan === "baller" && sub) {
    return {
      key: sub,
      requestsAllowed: BALLER_REQUESTS_PER_MINUTE,
      timeWindowMinutes: 1,
    };
  }

  return {
    key: sub ?? "anonymous",
    requestsAllowed: FREE_REQUESTS_PER_MINUTE,
    timeWindowMinutes: 1,
  };
}
