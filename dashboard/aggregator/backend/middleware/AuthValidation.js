export const signupValidation = (req, res, next) => {
  const { name, companyName, companyWebsite, companyEmail, role, password } =
    req.body;

  if (!name || !companyName || !companyWebsite || !companyEmail || !role || !password) {
    return res.json({ success: false, message: "All fields are required" });
  }

  next();
};

export const loginValidation = (req, res, next) => {
  const { companyEmail, password } = req.body;

  if (!companyEmail || !password) {
    return res.json({ success: false, message: "Company email and password required" });
  }

  next();
};
