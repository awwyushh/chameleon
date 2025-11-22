import mongoose from "mongoose";

const AttackSchema = new mongoose.Schema({
  time: { type: Date, default: Date.now },
  srcIp: String,
  rawInput: String,
  geo: {
    country: String,
    city: String,
  },
  classification: {
    label: String,
    confidence: Number,
  },
  deception: String,
  tarpitDelayMs: Number,
});

export default mongoose.models.Attack || mongoose.model("Attack", AttackSchema);
