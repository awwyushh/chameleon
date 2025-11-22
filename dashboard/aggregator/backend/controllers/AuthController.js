import { UserModel } from "../models/user.js";
import bcrypt from "bcryptjs";
import jwt from "jsonwebtoken";

export const signup = async (req, res) => {
  try {
    const { name, companyName, companyWebsite, companyEmail, role, password } =
      req.body;

    const exists = await UserModel.findOne({ companyEmail });
    if (exists) {
      return res.json({
        success: false,
        message: "Company email already registered",
      });
    }

    const hashedPass = await bcrypt.hash(password, 10);

    const user = new UserModel({
      name,
      companyName,
      companyWebsite,
      companyEmail,
      role,
      password: hashedPass,
    });

    await user.save();

    res.json({
      success: true,
      message: "Signup successful",
    });
  } catch (error) {
    console.log(error);
    res.json({ success: false, message: "Server error" });
  }
};

export const login = async (req, res) => {
  try {
    const { companyEmail, password } = req.body;

    const user = await UserModel.findOne({ companyEmail });
    if (!user) {
      return res.json({ success: false, message: "Invalid email or password" });
    }

    const valid = await bcrypt.compare(password, user.password);
    if (!valid) {
      return res.json({ success: false, message: "Invalid password" });
    }

    const token = jwt.sign(
      { id: user._id, email: user.companyEmail },
      process.env.JWT_SECRET,
      { expiresIn: "7d" }
    );

    res.json({
      success: true,
      jwtToken: token,
      email: user.companyEmail,
    });
  } catch (error) {
    console.log(error);
    res.json({ success: false, message: "Server error" });
  }
};
