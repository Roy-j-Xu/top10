

export const TopTenMsgType = {
  BROADCAST      : "topten:broadcast",
  START          : "topten:start",
  GAME_INFO      : "topten:game-info",
  START_GUESSING   : "topten:start-guessing",
  REVEAL_NUMBER  : "topten:reveal-number",
  FINISHED       : "topten:finished",
  ERROR          : "topten:error",
} as const

export const TopTenPlayerMsgType = {
  READY        : "topten-player:ready",
  SET_QUESTION : "topten-player:set-question",
  CHOOSE_ORDER : "topten-player:choose-order",
}

export interface GameInfo {
  turn: number;
  maxTurn: number;
  turnOrder: string[];
  guesser: string;
  questions: string[];
  usedQuestion: string;
  numbers: Record<string, number>;
  state: string;
}