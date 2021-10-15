{-# LANGUAGE DeriveGeneric #-}
{-# LANGUAGE FlexibleContexts #-}

module Main where

import Control.Monad.Writer
import Data.Aeson
import qualified Data.ByteString.Lazy as BSL
import GHC.Generics
import System.Console.ANSI
import System.Directory
import System.Exit
import System.IO
import System.Process

autogradingPath :: FilePath
autogradingPath = ".github/classroom/autograding.json"

data Test = Test
  { name :: String,
    setup :: String,
    run :: String,
    input :: String,
    output :: String,
    comparison :: String,
    timeout :: Integer,
    points :: Integer
  }
  deriving (Generic, Show)

instance FromJSON Test

newtype Tests = Tests
  { tests :: [Test]
  }
  deriving (Generic, Show)

instance FromJSON Tests

main :: IO ()
main = do
  exists <- doesFileExist autogradingPath
  if not exists
    then do
      putStrLn "tester: could not find file .github/classroom/autograding.json"
      exitFailure
    else do
      autograding <- BSL.readFile autogradingPath
      case decode autograding of
        Nothing -> do
          putStrLn "tester: could not decode test json data"
          exitFailure
        Just ts -> do
          (s, log) <- runWriterT $ runTests ts 0
          setSGR [Reset]
          putStr "\n\n"
          setSGR [SetColor Background Vivid Black]
          setSGR [SetColor Foreground Vivid White]
          putStr $ "Tests complete. Score: " ++ show s ++ "/100. Test log written to testlog.txt."
          setSGR [Reset]
          putStr "\n\n"
          writeFile "testlog.txt" (unlines log)

runTests :: Tests -> Integer -> WriterT [String] IO Integer
runTests (Tests []) s = return s
runTests (Tests (t : ts)) s =
  do
    logInfo $ "running: test " ++ name t ++ " with command '" ++ run t ++ "' ..."
    (_, Just hout, _, _) <- liftIO $ createProcess (shell (run t)) {std_out = CreatePipe}
    out <- liftIO $ hGetContents hout
    if out /= output t
      then do
        logFailure ("test " ++ name t ++ " failed. Output of command " ++ run t ++ " did not match expected. Score: " ++ (show s ++ "/100")) (output t) out
        runTests (Tests ts) s
      else do
        logSuccess $ "test " ++ name t ++ " succeeded! Score: " ++ (show (s + points t) ++ "/100")
        runTests (Tests ts) (s + points t)

logInfo s = do
  liftIO $ setSGR [Reset]
  liftIO $ putStrLn $ "tester: " ++ s
  tell ["tester: " ++ s]

logFailure s e o = do
  liftIO $ do
    setSGR [Reset]
    putStr "tester: "
    setSGR [SetColor Foreground Vivid Red]
    putStr "FAILURE"
    setSGR [Reset]
    putStrLn $ ": " ++ s
    putStrLn "\t\tExpected output:"
    putStr $ unlines $ map ('\t' :) (lines e)
    putStrLn "\t\tActual output:"
    putStr $ unlines $ map ('\t' :) (lines o)
  tell ["tester: FAILURE: " ++ s ++ "\n\t\tExpected output:\n" ++ unlines (map ('\t' :) (lines e)) ++ "\t\tActual output:\n" ++ (\s -> take (length s - 1) s) (unlines (map ('\t' :) (lines o)))]

logSuccess s = do
  liftIO $ do
    setSGR [Reset]
    putStr "tester: "
    setSGR [SetColor Foreground Dull Green]
    putStr "SUCCESS"
    setSGR [Reset]
    putStrLn $ ": " ++ s
  tell ["tester: SUCCESS: " ++ s]
