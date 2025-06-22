// MongoDB initialization script
// This script runs when the MongoDB container starts for the first time

// Switch to the hub_service database
db = db.getSiblingDB('hub_service');

// Create collections with validation
db.createCollection('users', {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["email", "password", "name", "role"],
            properties: {
                email: {
                    bsonType: "string",
                    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
                },
                password: {
                    bsonType: "string",
                    minLength: 6
                },
                name: {
                    bsonType: "string",
                    minLength: 1
                },
                role: {
                    enum: ["admin", "super_admin", "client"]
                }
            }
        }
    }
});

db.createCollection('challenges', {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["title", "content", "difficulty"],
            properties: {
                title: {
                    bsonType: "string",
                    minLength: 1
                },
                content: {
                    bsonType: "string",
                    minLength: 1
                },
                difficulty: {
                    enum: ["easy", "medium", "hard"]
                }
            }
        }
    }
});

// Create indexes for better performance
db.users.createIndex({ "email": 1 }, { unique: true });
db.users.createIndex({ "role": 1 });
db.challenges.createIndex({ "title": 1 });
db.challenges.createIndex({ "difficulty": 1 });

// Insert a default super admin user if not exists
const defaultAdmin = {
    email: "admin@hubservice.com",
    password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // "password" hashed with bcrypt
    name: "Super Admin",
    role: "super_admin",
    avatar: "",
    created_at: new Date(),
    updated_at: new Date()
};

// Check if admin user already exists
const existingAdmin = db.users.findOne({ email: defaultAdmin.email });
if (!existingAdmin) {
    db.users.insertOne(defaultAdmin);
    print("Default super admin user created: admin@hubservice.com / password");
} else {
    print("Super admin user already exists");
}

print("MongoDB initialization completed successfully!"); 