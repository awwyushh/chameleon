import mongoose from "mongoose"

export const connectDB = async() =>{
  try{
    await mongoose.connect(process.env.MONGO_URL);
    console.log("successfully connected DB!")
  }catch(error){
    console.log("Error detected",error)
  }
}
