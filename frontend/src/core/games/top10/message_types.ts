

export const TopTenMsgType = {
  BROADCAST      : "topten:broadcast",
  NEW_QUESTIONS  : "topten:new-questions",
  SET_QUESTION   : "topten:set-question",
  ASSIGN_NUMBERS : "topten:assign-numbers",
  FINISHED       : "topten:finished",
  ERROR          : "topten:error",

  P_READY        : "topten-player:ready",
  P_SET_QUESTION : "topten-player:set-question",
  P_CHOOSE_ORDER : "topten-player:choose-order",
} as const