import mongoose from "mongoose";

const Schema = mongoose.Schema;

const UserSchema = new Schema({
  name: String,
  companyName: String,
  companyWebsite: String,
  companyEmail: { type: String, unique: true },
  role: String,
  password: String,
});

// FIX: prevents OverwriteModelError
export const UserModel =
  mongoose.models.users || mongoose.model("users", UserSchema);
