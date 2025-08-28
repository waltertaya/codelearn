#!/usr/bin/env python3
"""
CodeLearn CLI - Submit your code solutions directly from your local environment
"""

import os
import sys
import json
import requests
import argparse
from pathlib import Path

# Configuration
API_BASE_URL = "http://localhost:8080/api/v1"
CONFIG_FILE = Path.home() / ".codelearn" / "config.json"

class CodeLearnCLI:
    def __init__(self):
        self.config = self.load_config()
        
    def load_config(self):
        """Load configuration from file"""
        if CONFIG_FILE.exists():
            with open(CONFIG_FILE, 'r') as f:
                return json.load(f)
        return {}
    
    def save_config(self, config):
        """Save configuration to file"""
        CONFIG_FILE.parent.mkdir(exist_ok=True)
        with open(CONFIG_FILE, 'w') as f:
            json.dump(config, f, indent=2)
        self.config = config
    
    def login(self, username, password):
        """Login and save authentication token"""
        try:
            response = requests.post(f"{API_BASE_URL}/auth/login", json={
                "username": username,
                "password": password
            })
            
            if response.status_code == 200:
                data = response.json()
                self.save_config({
                    "token": data["token"],
                    "username": data["user"]["username"],
                    "user_id": data["user"]["id"]
                })
                print(f"✅ Successfully logged in as {username}")
                return True
            else:
                print(f"❌ Login failed: {response.json().get('error', 'Unknown error')}")
                return False
        except Exception as e:
            print(f"❌ Error during login: {e}")
            return False
    
    def register(self, username, email, password):
        """Register a new user"""
        try:
            response = requests.post(f"{API_BASE_URL}/auth/register", json={
                "username": username,
                "email": email,
                "password": password
            })
            
            if response.status_code == 201:
                data = response.json()
                self.save_config({
                    "token": data["token"],
                    "username": data["user"]["username"],
                    "user_id": data["user"]["id"]
                })
                print(f"✅ Successfully registered and logged in as {username}")
                return True
            else:
                print(f"❌ Registration failed: {response.json().get('error', 'Unknown error')}")
                return False
        except Exception as e:
            print(f"❌ Error during registration: {e}")
            return False
    
    def get_headers(self):
        """Get authentication headers"""
        if "token" not in self.config:
            print("❌ Not logged in. Please run 'codelearn login' first.")
            return None
        return {"Authorization": f"Bearer {self.config['token']}"}
    
    def list_challenges(self, difficulty=None, language=None):
        """List available challenges"""
        headers = self.get_headers()
        if not headers:
            return
        
        params = {}
        if difficulty:
            params["difficulty"] = difficulty
        if language:
            params["language"] = language
        
        try:
            response = requests.get(f"{API_BASE_URL}/challenges", headers=headers, params=params)
            
            if response.status_code == 200:
                data = response.json()
                challenges = data["challenges"]
                
                print(f"\n📚 Available Challenges ({len(challenges)} found):")
                print("-" * 80)
                
                for challenge in challenges:
                    print(f"ID: {challenge['id']}")
                    print(f"Title: {challenge['title']}")
                    print(f"Difficulty: {challenge['difficulty']}")
                    print(f"Language: {challenge['language']}")
                    print(f"Description: {challenge['description'][:100]}...")
                    print("-" * 80)
            else:
                print(f"❌ Failed to fetch challenges: {response.json().get('error', 'Unknown error')}")
        except Exception as e:
            print(f"❌ Error fetching challenges: {e}")
    
    def get_challenge(self, challenge_id):
        """Get detailed challenge information"""
        headers = self.get_headers()
        if not headers:
            return None
        
        try:
            response = requests.get(f"{API_BASE_URL}/challenges/{challenge_id}", headers=headers)
            
            if response.status_code == 200:
                return response.json()
            else:
                print(f"❌ Failed to fetch challenge: {response.json().get('error', 'Unknown error')}")
                return None
        except Exception as e:
            print(f"❌ Error fetching challenge: {e}")
            return None
    
    def submit_solution(self, challenge_id, file_path, language=None):
        """Submit a solution file"""
        headers = self.get_headers()
        if not headers:
            return
        
        # Read the code file
        try:
            with open(file_path, 'r') as f:
                code = f.read()
        except FileNotFoundError:
            print(f"❌ File not found: {file_path}")
            return
        except Exception as e:
            print(f"❌ Error reading file: {e}")
            return
        
        # Detect language from file extension if not provided
        if not language:
            ext = Path(file_path).suffix.lower()
            language_map = {
                '.py': 'python',
                '.js': 'javascript',
                '.go': 'go',
                '.cs': 'csharp',
                '.c': 'c',
                '.cpp': 'cpp',
                '.rs': 'rust'
            }
            language = language_map.get(ext, 'python')
        
        # Get challenge details first
        challenge = self.get_challenge(challenge_id)
        if not challenge:
            return
        
        print(f"\n🚀 Submitting solution for: {challenge['title']}")
        print(f"Language: {language}")
        print(f"File: {file_path}")
        
        try:
            response = requests.post(f"{API_BASE_URL}/challenges/{challenge_id}/submit", 
                                   headers=headers, 
                                   json={
                                       "code": code,
                                       "language": language
                                   })
            
            if response.status_code == 201:
                data = response.json()
                submission = data["submission"]
                
                print(f"\n✅ Solution submitted successfully!")
                print(f"Submission ID: {submission['id']}")
                print(f"Status: {submission['status']}")
                print(f"Score: {submission['score']}/100")
                print(f"Output: {submission['output']}")
            else:
                print(f"❌ Submission failed: {response.json().get('error', 'Unknown error')}")
        except Exception as e:
            print(f"❌ Error submitting solution: {e}")
    
    def list_submissions(self):
        """List user's submissions"""
        headers = self.get_headers()
        if not headers:
            return
        
        try:
            response = requests.get(f"{API_BASE_URL}/submissions", headers=headers)
            
            if response.status_code == 200:
                data = response.json()
                submissions = data["submissions"]
                
                print(f"\n📝 Your Submissions ({len(submissions)} found):")
                print("-" * 80)
                
                for submission in submissions:
                    print(f"ID: {submission['id']}")
                    print(f"Challenge: {submission['challenge_title']}")
                    print(f"Language: {submission['language']}")
                    print(f"Status: {submission['status']}")
                    print(f"Score: {submission['score']}/100")
                    print(f"Submitted: {submission['created_at']}")
                    print("-" * 80)
            else:
                print(f"❌ Failed to fetch submissions: {response.json().get('error', 'Unknown error')}")
        except Exception as e:
            print(f"❌ Error fetching submissions: {e}")
    
    def show_leaderboard(self):
        """Show the leaderboard"""
        headers = self.get_headers()
        if not headers:
            return
        
        try:
            response = requests.get(f"{API_BASE_URL}/leaderboard", headers=headers)
            
            if response.status_code == 200:
                data = response.json()
                leaderboard = data["leaderboard"]
                
                print(f"\n🏆 Leaderboard (Top {len(leaderboard)}):")
                print("-" * 80)
                print(f"{'Rank':<6} {'Username':<20} {'Score':<10} {'Submissions':<12} {'Last Activity'}")
                print("-" * 80)
                
                for i, entry in enumerate(leaderboard, 1):
                    print(f"{i:<6} {entry['username']:<20} {entry['total_score']:<10} {entry['submissions']:<12} {entry['last_activity']}")
            else:
                print(f"❌ Failed to fetch leaderboard: {response.json().get('error', 'Unknown error')}")
        except Exception as e:
            print(f"❌ Error fetching leaderboard: {e}")

def main():
    parser = argparse.ArgumentParser(description="CodeLearn CLI - Code locally, learn globally")
    subparsers = parser.add_subparsers(dest="command", help="Available commands")
    
    # Login command
    login_parser = subparsers.add_parser("login", help="Login to your account")
    login_parser.add_argument("username", help="Your username")
    login_parser.add_argument("password", help="Your password")
    
    # Register command
    register_parser = subparsers.add_parser("register", help="Register a new account")
    register_parser.add_argument("username", help="Choose a username")
    register_parser.add_argument("email", help="Your email address")
    register_parser.add_argument("password", help="Choose a password")
    
    # List challenges command
    list_parser = subparsers.add_parser("challenges", help="List available challenges")
    list_parser.add_argument("--difficulty", choices=["Easy", "Medium", "Hard"], help="Filter by difficulty")
    list_parser.add_argument("--language", help="Filter by programming language")
    
    # Submit command
    submit_parser = subparsers.add_parser("submit", help="Submit a solution")
    submit_parser.add_argument("challenge_id", type=int, help="Challenge ID")
    submit_parser.add_argument("file", help="Path to your solution file")
    submit_parser.add_argument("--language", help="Programming language (auto-detected if not specified)")
    
    # Submissions command
    subparsers.add_parser("submissions", help="List your submissions")
    
    # Leaderboard command
    subparsers.add_parser("leaderboard", help="Show the leaderboard")
    
    args = parser.parse_args()
    
    if not args.command:
        parser.print_help()
        return
    
    cli = CodeLearnCLI()
    
    if args.command == "login":
        cli.login(args.username, args.password)
    elif args.command == "register":
        cli.register(args.username, args.email, args.password)
    elif args.command == "challenges":
        cli.list_challenges(args.difficulty, args.language)
    elif args.command == "submit":
        cli.submit_solution(args.challenge_id, args.file, args.language)
    elif args.command == "submissions":
        cli.list_submissions()
    elif args.command == "leaderboard":
        cli.show_leaderboard()

if __name__ == "__main__":
    main()
