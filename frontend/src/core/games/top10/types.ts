

export const TopTenMsgType = {
  BROADCAST      : "topten:broadcast",
  TURN_INFO      : "topten:turn-info",
  SET_QUESTION   : "topten:set-question",
  ASSIGN_NUMBERS : "topten:assign-numbers",
  REVEAL_NUMBER  : "topten:reveal-number",
  FINISHED       : "topten:finished",
  ERROR          : "topten:error",
} as const

export const TopTenPlayerMsgType = {
  READY        : "topten-player:ready",
  SET_QUESTION : "topten-player:set-question",
  CHOOSE_ORDER : "topten-player:choose-order",
}

export interface TurnInfoMsgData {
  turn: number;
  guesser: string;
  questions: string[];
}